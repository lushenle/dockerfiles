package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

// Config 导出sql所需的配置信息
type Config struct {
	Debug          bool      // 是否调试模式
	IsExportData   bool      // 是否导出数据
	ExportDataStep int64     // 导出数据时，每次查询数据量
	IsCreateDB     bool      // 是否生成建库语句
	OutPath        string    // 输出sql文件目录-绝对路径-用于导出
	OutZip         bool      // 是否导出zip压缩文件
	DbCfg          *DbConfig // 数据库连接信息
}

// DbConfig 数据库连接配置
type DbConfig struct {
	Address string // 数据库连接地址
	Port    int    // 数据库端口
	User    string // 数据库用户名
	Passwd  string // 数据库密码
	DbName  string // 数据库名
}

// CreateTable 创建table sql查询
type CreateTable struct {
	Table       string `xorm:"'Table'"`
	CreateTable string `xorm:"'Create Table'"`
}

// CreateDb 创建数据库 sql查询
type CreateDb struct {
	Database       string `xorm:"'Database'"`
	CreateDatabase string `xorm:"'Create Database'"`
}

// TPLSqlModel 导出数据sql部分
type TPLSqlModel struct {
	TableName string // 表名
	CreateSQL string // 创建表sql语句
	InsertSQL string // 插入数据sql语句
}

// TPLModel 导出数据结构体
type TPLModel struct {
	MySQL *DbConfig
	SQL   []*TPLSqlModel
	Date  string
}

// TableColumn 用于读取数据库每一列数据
type TableColumn interface {
}

// Mysqldump mysql导出数据对象
type Mysqldump struct {
	conn    *xorm.Engine
	cfg     *Config
	isClose bool
}

// New 创建一个Mysqldump对象
func New(cfg *Config) (*Mysqldump, error) {
	if cfg == nil {
		return nil, errors.New("配置信息不能为nil")
	}
	// 处理配置信息
	if cfg.OutPath == "" {
		return nil, errors.New("导出sql输出路径不能是空")
	}
	if cfg.ExportDataStep == 0 {
		cfg.ExportDataStep = 1000
	}

	// 创建导出对象
	mysqldump := &Mysqldump{
		cfg:     cfg,
		isClose: false,
	}
	// 连接mysql
	err := mysqldump.OpenMysql()
	if err != nil {
		return nil, err
	}

	return mysqldump, nil
}

// Close 不使用导出功能时，关闭连接资源
func (md *Mysqldump) Close() error {
	md.isClose = true
	return md.conn.Close()
}

// OpenMysql 连接mysql
func (md *Mysqldump) OpenMysql() error {
	// 拼接连接数据库字符串
	connStr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=UTC",
		md.cfg.DbCfg.User,
		md.cfg.DbCfg.Passwd,
		md.cfg.DbCfg.Address,
		md.cfg.DbCfg.Port,
		md.cfg.DbCfg.DbName)

	// 连接数据库
	engine, err := xorm.NewEngine("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// 是否开启debug模式
	if md.cfg.Debug {
		engine.Logger().SetLevel(core.LOG_DEBUG) // 调试信息
		engine.ShowSQL(true)                     // 显示sql
	}
	engine.SetMaxIdleConns(2)            // 空闲连接池数量
	engine.SetMaxOpenConns(8)            // 最大连接数
	engine.SetMapper(core.GonicMapper{}) // 命名规则

	// 设置数据库时区
	engine.DatabaseTZ = time.UTC
	engine.TZLocation = time.UTC

	md.conn = engine

	log.Println("连接数据库成功")
	return nil
}

// GetRootDir 获取程序跟目录,返回值尾部包含'/'
func (md *Mysqldump) GetRootDir() string {
	// 文件不存在获取执行路径
	file, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		file = fmt.Sprintf(".%s", string(os.PathSeparator))
	} else {
		file = fmt.Sprintf("%s%s", file, string(os.PathSeparator))
	}
	return file
}

// Export 导出数据库所有表
func (md *Mysqldump) Export() (outFile string, err error) {
	if md.isClose == true {
		return "", errors.New("已调用Close关闭相关资源，无法进行导出")
	}
	// 创建导出sql文件
	outFile = fmt.Sprintf("%s/%s_%s.sql", strings.TrimRight(md.cfg.OutPath, "/"), md.cfg.DbCfg.DbName, time.Now().Format("2006-01-02"))
	lf, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal(err)
	}

	// 获取建库语句
	createSQL, err := md.GetCreateDbSQL()
	if err != nil {
		log.Fatal(err)
	}
	// 获取数据库字符集
	charSet := "utf8"
	var valid = regexp.MustCompile("CHARACTER SET ([a-z0-9A-Z]+) ")
	finds := valid.FindAllStringSubmatch(createSQL, -1)
	if len(finds) > 0 {
		if len(finds[0]) > 1 {
			charSet = finds[0][1]
		}
	}
	// 判断是否需要建库语句
	if md.cfg.IsCreateDB == true {
		createSQL = fmt.Sprintf("CREATE DATABASE `%s`; /*!40100 DEFAULT CHARACTER SET %s */", md.cfg.DbCfg.DbName, charSet)
	} else {
		createSQL = ""
	}

	// 写入头部信息
	_, err = lf.WriteString(fmt.Sprintf(`/*
		databaseAddress        : %s:%d
		database         : %s
		createTime: %s
*/
%s
SET NAMES %s;
SET FOREIGN_KEY_CHECKS = 0;
`,
		md.cfg.DbCfg.Address,
		md.cfg.DbCfg.Port,
		md.cfg.DbCfg.DbName,
		time.Now().Format("2006-01-02"),
		createSQL,
		charSet))
	if err != nil {
		return "", err
	}

	// 查询数据库表列表
	tables, err := md.SelectTableNames()
	if err != nil {
		log.Fatal(err)
	}

	// 导出数据对象
	tplSQLModel := make([]*TPLSqlModel, 0)
	// 循环表名，查询出对应的表创建语句
	for _, table := range tables {
		log.Println("dump: ", table)
		// 导出建表语句
		sql, err := md.GetCreateTableSQL(table)
		if err != nil {
			log.Fatal(err)
		}
		tplSQLModel = append(tplSQLModel, &TPLSqlModel{
			TableName: table,
			CreateSQL: sql,
		})
		log.Println(sql)
		// 写入一个表到建表语句
		_, err = lf.WriteString(fmt.Sprintf(
			`%s-- ----------------------------
-- Table structure for %s
-- ----------------------------
DROP TABLE IF EXISTS %s%s%s;
%s;
%s`,
			"\n\n",
			table,
			"`",
			table,
			"`",
			sql,
			"\n"))
		if err != nil {
			log.Fatal(err)
		}
		// 导出数据
		if md.cfg.IsExportData == true {
			md.ExportData(lf, table)
		}
	}

	return outFile, nil
}

// SelectTableNames 查询数据库表列表
func (md *Mysqldump) SelectTableNames() (tables []string, err error) {
	tables = make([]string, 0)
	err = md.conn.SQL("SHOW TABLES;").Cols(fmt.Sprintf("Tables_in_%s", md.cfg.DbCfg.DbName)).Find(&tables)
	return
}

// GetCreateTableSQL 查询创建表语句
func (md *Mysqldump) GetCreateTableSQL(tableName string) (string, error) {
	creates := make([]*CreateTable, 0)
	err := md.conn.SQL(fmt.Sprintf("show create table %s", tableName)).Find(&creates)
	log.Println(err)
	if err != nil {
		return "", err
	}
	if len(creates) == 0 {
		return "", errors.New("查询table 创建语句错误")
	}
	return creates[0].CreateTable, nil
}

// GetCreateDbSQL 获取创建数据库
func (md *Mysqldump) GetCreateDbSQL() (string, error) {
	createSQLs := make([]*CreateDb, 0)
	err := md.conn.SQL(fmt.Sprintf("SHOW CREATE DATABASE %s", md.cfg.DbCfg.DbName)).Find(&createSQLs)
	if err != nil {
		return "", err
	}
	if len(createSQLs) == 0 {
		return "", errors.New("查询创建数据库语句为空")
	}
	return createSQLs[0].CreateDatabase, nil
}

// ExportData 导出数据为
func (md *Mysqldump) ExportData(w io.Writer, tableName string) (err error) {
	log.Println("Starting dump table:", tableName)
	// 查询总数据行数
	var count int64
	count, err = md.conn.Table(tableName).Count()
	if err != nil {
		return
	}
	log.Println(count)

	columns, xormColumns, err := md.conn.Dialect().GetColumns(tableName)
	if err != nil {
		return err
	}

	var offset int64
	for offset = 0; offset < count; offset += md.cfg.ExportDataStep {
		colNames := md.conn.Dialect().Quote(strings.Join(columns, md.conn.Dialect().Quote(", ")))
		sql := fmt.Sprintf("select %s from %s limit %d offset %d", colNames, tableName, md.cfg.ExportDataStep, offset)
		list, err := md.conn.QueryInterface(sql)
		if err != nil {
			return err
		}
		for _, one := range list {
			// 拼接插入语句头部
			installSQL := fmt.Sprintf("\nINSERT INTO %s (%s) VALUES ",
				md.conn.Dialect().Quote(tableName),
				colNames)
			values := make([]string, 0)
			for _, column := range columns {
				val, ok := one[column] // 读取本行值
				if ok == false {
					return errors.New("列名和值无法对应")
				}
				// 判断是否是时间类型
				if xormColumn, ok := xormColumns[column]; ok == true {
					aa[xormColumn.SQLType.Name] = xormColumn.SQLType.Name
					if xormColumn.SQLType.IsTime() == true {
						isTimeNull := false
						if val == nil {
							val = "null"
							isTimeNull = true
						} else {
							valTime := val.(time.Time)
							if err == nil {
								val = valTime.Format("2006-01-02")
							} else {
								val = "null"
								isTimeNull = true
							}
						}
						if isTimeNull == true {
							values = append(values, fmt.Sprintf("%v", val))
						} else {
							values = append(values, fmt.Sprintf("'%v'", val))
						}

					} else if xormColumn.SQLType.IsBlob() == true {
						if val == nil {
							val = "false"
						} else {
							if reflect.TypeOf(val).Kind() == reflect.Slice {
								val = md.conn.Dialect().FormatBytes(val.([]byte))
							} else if reflect.TypeOf(val).Kind() == reflect.String {
								val = val.(string)
							}
						}

						values = append(values, fmt.Sprintf("%v", val))
					} else if xormColumn.SQLType.IsNumeric() == true {
						if val == nil {
							val = "null"
						} else {
							if valByte, ok := val.([]byte); ok == true {
								// log.Println(column, "3-1")
								val = string(valByte)
							} else {
								// log.Println(column, "3-2")
								val = fmt.Sprint(val)
							}
						}

						values = append(values, fmt.Sprintf("%v", val))
					} else {
						if val == nil {
							val = ""
						} else {
							if valByte, ok := val.([]byte); ok == true {
								// log.Println(column, "3-1")
								val = string(valByte)
							} else {
								// log.Println(column, "3-2")
								val = fmt.Sprint(val)
							}
						}

						values = append(values, fmt.Sprintf("'%v'", val))
					}
				}
			}
			// 拼接插入语句值部分
			installSQL = fmt.Sprintf("%s (%s);", installSQL, strings.Join(values, ","))
			// 写入数据
			_, err = io.WriteString(w, installSQL)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var aa map[string]string

func init() {
	aa = make(map[string]string, 0)
}

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())

	// 系统日志显示文件和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	cfg := &Config{
		Debug:        false,
		IsExportData: true,
		IsCreateDB:   false,
		OutZip:       false,
		OutPath:      "/mnt/data1/oa-erp/oa/mysql",
		DbCfg: &DbConfig{
			Address: "x.x.x.x",
			Port:    3306,
			User:    "root",
			Passwd:  "xxxx",
			DbName:  "xxx",
		},
	}
	dm, err := New(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	// 导出
	path, err := dm.Export()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(path)
	//select {}
}
