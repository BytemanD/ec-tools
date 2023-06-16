module github.com/BytemanD/ec-tools

go 1.17

require (
	github.com/spf13/cobra v1.7.0
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jedib0t/go-pretty/v6 v6.4.6 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	libvirt.org/go/libvirt v1.9002.0 // indirect
)

require (
	github.com/BytemanD/ec-tools/pkg v0.0.0
	github.com/BytemanD/easygo/pkg v0.0.3
)

replace github.com/BytemanD/ec-tools/pkg => ./pkg
