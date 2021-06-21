EXE = ghfs
SRC = .
LDFLAGS = -ldflags="-s -w"

define embed
    sed -i -e '/= `/,/`/ { /`/!d }' $(1).go
    sed -i -e '/= `/r$(1)' $(1).go
    sed -i -e '/= `/,/`/ { /`/! { /^\s*$$/d; s/^\s\s*// } }' $(1).go
endef

tpls:
    $(call embed, tpl/page.html)
    sed -i -e 's;/../asset/;?asset=;g' tpl/page.html.go

    $(call embed, tpl/asset/main.css)

    $(call embed, tpl/asset/main.js)


windows:
    GO111MODULE=off GOOS=windows go build -o $(EXE)_win.exe $(LDFLAGS) $(SRC)

macos:
    GO111MODULE=off GOOS=darwin go build -o $(EXE)_macos $(LDFLAGS) $(SRC)

linux:
    GO111MODULE=off GOOS=linux go build -o $(EXE)_linux $(LDFLAGS) $(SRC)

all: windows macos linux
    echo "done."

clean:
    rm -f $(EXE)_win.exe $(EXE)_macos $(EXE)_linux

