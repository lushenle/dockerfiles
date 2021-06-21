package main

import (
	"bufio"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	currentUser *user.User
	hostInfo    = make(map[string]Section)
)

// Section ssh 相关参数
type Section struct {
	Hostname     string
	Port         int
	User         string
	IdentityFile string
}

func (s *Section) clear() {
	s.Hostname = ""
	s.Port = 0
	s.User = ""
	s.IdentityFile = ""
}

func dropError(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func getpass(prompt string) (pass string, err error) {

	tstate, err := terminal.GetState(0)
	dropError(err)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		quit := false
		for _ = range sig {
			quit = true
			break
		}
		terminal.Restore(0, tstate)
		if quit {
			fmt.Println()
			os.Exit(2)
		}
	}()
	defer func() {
		signal.Stop(sig)
		close(sig)
	}()

	f := bufio.NewWriter(os.Stdout)
	f.Write([]byte(prompt))
	f.Flush()

	passbytes, err := terminal.ReadPassword(0)
	dropError(err)
	pass = string(passbytes)

	f.Write([]byte("\n"))
	f.Flush()

	return
}

func parsePemBlock(block *pem.Block) (interface{}, error) {

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "DSA PRIVATE KEY":
		return ssh.ParseDSAPrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("gssh: unsupported key type %q", block.Type)
	}
}

func expandPath(path string) string {

	if len(path) < 2 || path[:2] != "~/" {
		return path
	}
	return strings.Replace(path, "~", currentUser.HomeDir, 1)
}

func addKeyAuth(auths []ssh.AuthMethod, keypath string) []ssh.AuthMethod {
	if len(keypath) == 0 {
		return auths
	}

	keypath = expandPath(keypath)

	// read the file
	pemBytes, err := ioutil.ReadFile(keypath)
	dropError(err)

	// get first pem block
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		log.Printf("no key found in %s", keypath)
		return auths
	}

	// handle plain and encrypted keyfiles
	if x509.IsEncryptedPEMBlock(block) {
		prompt := fmt.Sprintf("Enter passphrase for key '%s': ", keypath)
		pass, err := getpass(prompt)
		dropError(err)

		block.Bytes, err = x509.DecryptPEMBlock(block, []byte(pass))
		dropError(err)

		key, err := parsePemBlock(block)
		dropError(err)

		signer, err := ssh.NewSignerFromKey(key)
		dropError(err)

		return append(auths, ssh.PublicKeys(signer))
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	dropError(err)
	return append(auths, ssh.PublicKeys(signer))
}

func getAgentAuth() (auth ssh.AuthMethod, ok bool) {
	if sock := os.Getenv("SSH_AUTH_SOCK"); len(sock) > 0 {
		if agconn, err := net.Dial("unix", sock); err == nil {
			ag := agent.NewClient(agconn)
			auth = ssh.PublicKeysCallback(ag.Signers)
			ok = true
		}
	}
	return
}

func addPasswordAuth(user, addr string, auths []ssh.AuthMethod) []ssh.AuthMethod {
	if terminal.IsTerminal(0) == false {
		return auths
	}
	host := addr
	if i := strings.LastIndex(host, ":"); i != -1 {
		host = host[:i]
	}
	prompt := fmt.Sprintf("%s@%s's password: ", user, host)
	passwordCallback := func() (string, error) {
		return getpass(prompt)
	}
	return append(auths, ssh.PasswordCallback(passwordCallback))
}

func tryAgentConnect(user, addr string) (client *ssh.Client) {
	if auth, ok := getAgentAuth(); ok {
		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{auth},
		}
		client, _ = ssh.Dial("tcp", addr, config)
	}
	return
}

func sshConnect(user, addr, keypath string) (client *ssh.Client) {
	// try connecting via agent first
	client = tryAgentConnect(user, addr)
	if client != nil {
		return
	}

	// if that failed try with the key and password methods
	auths := make([]ssh.AuthMethod, 0, 2)
	auths = addKeyAuth(auths, keypath)
	auths = addPasswordAuth(user, addr, auths)

	config := &ssh.ClientConfig{
		User: user,
		Auth: auths,
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", addr, config)
	dropError(err)

	return
}

func (s *Section) getFull(name string, def Section) (host string, port int, user, keyfile string) {
	if len(s.Hostname) > 0 {
		host = s.Hostname
	} else if len(def.Hostname) > 0 {
		host = def.Hostname
	}
	if s.Port > 0 {
		port = s.Port
	} else if def.Port > 0 {
		port = def.Port
	}
	if len(s.User) > 0 {
		user = s.User
	} else if len(def.User) > 0 {
		user = def.User
	}
	if len(s.IdentityFile) > 0 {
		keyfile = s.IdentityFile
	} else if len(def.IdentityFile) > 0 {
		keyfile = def.IdentityFile
	}
	return
}

func getSSHEntry(name string) (host string, port int, user, keyfile string) {

	def := Section{Hostname: name}
	if defcfg, ok := hostInfo["*"]; ok {
		def = defcfg
	}

	if s, ok := hostInfo[name]; ok {
		return s.getFull(name, def)
	}
	for h, s := range hostInfo {
		if ok, err := path.Match(h, name); ok && err == nil {
			return s.getFull(name, def)
		}
	}
	return def.Hostname, def.Port, def.User, def.IdentityFile
}

func parseSSHConfig(path string) bool {
	f, err := os.Open(path)
	dropError(err)
	defer f.Close()
	update := func(cb func(s *Section)) {}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 1 && strings.ToLower(parts[0]) == "host" {
			hosts := parts[1:]
			for _, h := range hosts {
				if _, ok := hostInfo[h]; !ok {
					hostInfo[h] = Section{}
				}
			}
			update = func(cb func(s *Section)) {
				for _, h := range hosts {
					s, _ := hostInfo[h]
					cb(&s)
					hostInfo[h] = s
				}
			}
		}
		if len(parts) == 2 {
			switch strings.ToLower(parts[0]) {
			case "hostname":
				update(func(s *Section) {
					s.Hostname = parts[1]
				})
			case "port":
				if p, err := strconv.Atoi(parts[1]); err == nil {
					update(func(s *Section) {
						s.Port = p
					})
				}
			case "user":
				update(func(s *Section) {
					s.User = parts[1]
				})
			case "identityfile":
				update(func(s *Section) {
					s.IdentityFile = parts[1]
				})
			}
		}
	}
	return true
}

func usage(code int) {
	fmt.Printf(
		`
Usage: gssh [-i private-key-file] [user@]host[:port]
	-i private-key-file
		PEM-encoded private key file to use (default: ~/.ssh/id_rsa if present)
	[user@]host[:port]
		the SSH server to connect to, with optional username and port
`)
	os.Exit(code)
}

func shift(q []string) (ok bool, val string, qnew []string) {
	if len(q) > 0 {
		ok = true
		val = q[0]
		qnew = q[1:]
	}
	return
}

func parseCmdLine() (host string, port int, user, key string) {
	ok, arg, args := shift(os.Args)
	var argKey, argHost, argInt string
	for ok {
		ok, arg, args = shift(args)
		if !ok {
			break
		}
		if arg == "-h" || arg == "--help" || arg == "--version" {
			usage(0)
		}
		if arg == "-i" {
			ok, argKey, args = shift(args)
			if !ok {
				usage(1)
			}
		} else if len(argHost) == 0 {
			argHost = arg
		} else if len(argInt) == 0 {
			argInt = arg
		} else {
			usage(1)
		}
	}
	if len(argHost) == 0 || argHost[0] == '-' {
		usage(1)
	}

	// key
	if len(argKey) != 0 {
		key = argKey
	} // else key remains ""

	// user, addr
	var addr string
	if i := strings.Index(argHost, "@"); i != -1 {
		user = argHost[:i]
		if i+1 >= len(argHost) {
			usage(1)
		}
		addr = argHost[i+1:]
	} else {
		addr = argHost
	}

	// addr -> host, port
	if p := strings.Split(addr, ":"); len(p) == 2 {
		host = p[0]
		var err error
		if port, err = strconv.Atoi(p[1]); err != nil {
			log.Printf("bad port: %v", err)
			usage(1)
		}
		if port <= 0 || port >= 65536 {
			log.Printf("bad port: %d", port)
			usage(1)
		}
	} else {
		host = addr
	}
	return
}

func main() {
	// get params from command line
	host, port, username, key := parseCmdLine()

	// get current user
	var err error
	currentUser, err = user.Current()
	dropError(err)

	sshConfig := filepath.Join(currentUser.HomeDir, ".ssh", "config")
	if _, err := os.Stat(sshConfig); err == nil {
		if parseSSHConfig(sshConfig) {
			shost, sport, suser, skey := getSSHEntry(host)
			if len(shost) > 0 {
				host = shost
			}
			if sport != 0 && port == 0 {
				port = sport
			}
			if len(suser) > 0 && len(username) == 0 {
				username = suser
			}
			if len(skey) > 0 && len(key) == 0 {
				key = skey
			}
		}
	}

	if port == 0 {
		port = 22
	}
	if len(username) == 0 {
		username = currentUser.Username
	}
	if len(key) == 0 {
		idrsap := filepath.Join(currentUser.HomeDir, ".ssh", "id_rsa")
		if _, err := os.Stat(idrsap); err == nil {
			key = idrsap
		}
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	client := sshConnect(username, addr, key)

	runCommand(client, "whoami")
}

// runCommand 远程执行命令的方法
func runCommand(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	dropError(err)
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	dropError(err)
	stdout = string(buf.Bytes())
	fmt.Println(stdout)
	return
}
