# Yggdrasil

[![CircleCI](https://circleci.com/gh/yggdrasil-network/yggdrasil-go.svg?style=shield&circle-token=:circle-token
)](https://circleci.com/gh/yggdrasil-network/yggdrasil-go)

## What is it?

This is a toy implementation of an encrypted IPv6 network, with many good ideas stolen from [cjdns](https://github.com/cjdelisle/cjdns), which was written to test a particular routing scheme that was cobbled together one random afternoon.
It's notably not a shortest path routing scheme, with the goal of scalable name-independent routing on dynamic networks with an internet-like topology.
It's named Yggdrasil after the world tree from Norse mythology, because that seemed like the obvious name given how it works.
More information is available at <https://yggdrasil-network.github.io/>.

This is a toy / proof-of-principle, and considered alpha quality by the developers. It's not expected to be feature complete, and future updates may not be backwards compatible, though it should warn you if it sees a connection attempt with a node running a newer version.
You're encouraged to play with it, but it is strongly advised not to use it for anything mission critical.

## Building

1. Install Go (tested on 1.11+, [godeb](https://github.com/niemeyer/godeb) is recommended for debian-based linux distributions).
2. Clone this repository.
2. `./build`

Note that you can cross-compile for other platforms and architectures by specifying the `$GOOS` and `$GOARCH` environment variables, for example, `GOOS=windows ./build` or `GOOS=linux GOARCH=mipsle ./build`.

The build script assumes you're using Go modules with Go 1.11+, so the build environment is self-contained when cloned outside of a GOPATH. You may build in a GOPATH, but this can cause yggdrasil to leak dependencies.

## Running

To run the program, you'll need permission to create a `tun` device and configure it using `ip`.
If you don't want to mess with capabilities for the `tun` device, then using `sudo` should work, with the usual security caveats about running a program as root.

To run with default settings:

1. `./yggdrasil --autoconf`

That will generate a new set of keys (and an IP address) each time the program is run.
The program will bind to all addresses on a random port and listen for incoming connections.
It will send announcements over IPv6 link-local multicast, and it will attempt to start a connection if it hears an announcement from another device.

In practice, you probably want to run this instead:

1. `./yggdrasil --genconf > conf.json`
2. `./yggdrasil --useconf < conf.json`

This keeps a persistent set of keys (and by extension, IP address) and gives you the option of editing the configuration file.
If you want to use it as an overlay network on top of e.g. the internet, then you can do so by adding the remote devices domain/address and port (as a string, e.g. `"1.2.3.4:5678"`) to the list of `Peers` in the configuration file.
By default, it peers over TCP (which can be forced with `"tcp://1.2.3.4:5678"` syntax), but it's also possible to connect over a socks proxy (`"socks://socksHost:socksPort/1.2.3.4:5678"`).
The socks proxy approach is useful for e.g. [peering over tor hidden services](https://github.com/yggdrasil-network/public-peers/blob/master/other/tor.md).
UDP support was removed as part of v0.2, and may be replaced by a better implementation at a later date.

### Platforms

#### Linux

- Should work out of the box on most Linux distributions with `iproute2` installed.
- systemd service scripts are included in the `contrib/systemd/` folder so that it runs automatically in the background (using `/etc/yggdrasil.conf` for configuration), copy the service files into `/etc/systemd/system`, copy `yggdrasil` into your `$PATH`, i.e. `/usr/bin`, and then enable the service:
```
systemctl enable yggdrasil
systemctl start yggdrasil
```
- Once installed as a systemd service, you can read the `yggdrasil` output:
```
systemctl status yggdrasil
journalctl -u yggdrasil
```

#### macOS

- Tested and working out of the box on macOS 10.13 High Sierra.
- May work in theory on any macOS version with `utun` support (which was added in macOS 10.7 Lion), although this is untested at present.
- TAP mode is not supported on macOS.

#### FreeBSD, NetBSD

- Works in TAP mode, but currently doesn't work in TUN mode.
- You may need to create the TAP adapter first if it doesn't already exist, i.e. `ifconfig tap0 create`.

#### OpenBSD

- Works in TAP mode, but currently doesn't work in TUN mode.
- You may need to create the TAP adapter first if it doesn't already exist, i.e. `ifconfig tap0 create`.
- OpenBSD is not capable of listening on both IPv4 and IPv6 at the same time on the same socket (unlike FreeBSD and NetBSD). This affects the `Listen` and `AdminListen` configuration options. You will need to set `Listen` and `AdminListen` to use either an IPv4 or an IPv6 address.
- You may consider using [relayd](https://man.openbsd.org/relayd.8) to allow incoming Yggdrasil connections on both IPv4 and IPv6 simultaneously.

#### Windows

- Tested and working on Windows 7 and Windows 10, and should work on any recent versions of Windows, but it depends on the [OpenVPN TAP driver](https://openvpn.net/index.php/open-source/downloads.html) being installed first.
- Has been proven to work with both the [NDIS 5](https://swupdate.openvpn.org/community/releases/tap-windows-9.9.2_3.exe) (`tap-windows-9.9.2_3`) driver and the [NDIS 6](https://swupdate.openvpn.org/community/releases/tap-windows-9.21.2.exe) (`tap-windows-9.21.2`) driver, however there are substantial performance issues with the NDIS 6 driver therefore it is recommended to use the NDIS 5 driver instead.
- Be aware that connectivity issues can occur on Windows if multiple IPv6 addresses from the `200::/7` prefix are assigned to the TAP interface. If this happens, then you may need to manually remove the old/unused addresses from the interface (though the code has a workaround in place to do this automatically in some cases).
- TUN mode is not supported on Windows.
- Yggdrasil can be installed as a Windows service so that it runs automatically in the background. From an Administrator Command Prompt:
```
sc create yggdrasil binpath= "\"C:\path\to\yggdrasil.exe\" -useconffile \"C:\path\to\yggdrasil.conf\""
sc config yggdrasil displayname= "Yggdrasil Service"
sc config yggdrasil start= "auto"
sc start yggdrasil
```
- Alternatively, if you want the service to autoconfigure instead of using an `yggdrasil.conf`, replace the `sc create` line from above with:
```
sc create yggdrasil binpath= "\"C:\path\to\yggdrasil.exe\" -autoconf"
```

#### EdgeRouter

- Tested and working on the EdgeRouter X, using the [vyatta-yggdrasil](https://github.com/neilalexander/vyatta-yggdrasil) wrapper package.

## Optional: advertise a prefix locally

Suppose a node has generated the address: `200:1111:2222:3333:4444:5555:6666:7777`

Then the node may also use addresses from the prefix: `300:1111:2222:3333::/64` (note the `200` changed to `300`, a separate `/8` is used for prefixes, but the rest of the first 64 bits are the same).

To advertise this prefix and a route to `200::/7`, the following seems to work on the developers' networks:

1. Enable IPv6 forwarding (e.g. `sysctl -w net.ipv6.conf.all.forwarding=1` or add it to sysctl.conf).

2. `ip addr add 300:1111:2222:3333::1/64 dev eth0` or similar, to assign an address for the router to use in that prefix, where the LAN is reachable through `eth0`.

3. Install/run `radvd` with something like the following in `/etc/radvd.conf`:
```
interface eth0
{
        AdvSendAdvert on;
        prefix 300:1111:2222:3333::/64 {
            AdvOnLink on;
            AdvAutonomous on;
        };
        route 200::/7 {};
};
```

This is enough to give unsupported devices on the LAN access to the yggdrasil network. See the [configuration](https://yggdrasil-network.github.io/configuration.html) page for more info.

## How does it work?

I'd rather not try to explain in the readme, but it is described further on the [about](https://yggdrasil-network.github.io/about.html) page, so you can check there if you're interested.
Be warned that it's still not a very good explanation, but it at least gives a high-level overview and links to some relevant work by other people.

## Obligatory performance propaganda

A [simplified model](misc/sim/treesim-forward.py) of this routing scheme has been tested in simulation on the 9204-node [skitter](https://www.caida.org/tools/measurement/skitter/) network topology dataset from [caida](https://www.caida.org/), and compared with results in [arxiv:0708.2309](https://arxiv.org/abs/0708.2309).
Using the routing scheme as implemented in this code, the average multiplicative stretch is observed to be about 1.08, with an average routing table size of 6 for a name-dependent scheme, and approximately 30 additional (but smaller) entries needed for the name-independent routing table.
The number of name-dependent routing table entries needed is proportional to node degree, so that 6 is the mean of a distribution with a long tail, but this may be an acceptable tradeoff(it's at least worth trying, hence this code).
The size of name-dependent routing table entries is relatively large, due to cryptographic signatures associated with routing table updates, but in the absence of cryptographic overhead, each entry should otherwise be comparable in size to the BC routing scheme described in the above paper.
A modified version of this scheme, with the same resource requirements, achieves a multiplicative stretch of 1.02, which drops to 1.01 if source routing is used.
Both of these optimizations are not present in the current implementation, as the former depends on network state information that appears difficult to cryptographically secure, and the latter optimization is both tedious to implement and would make debugging other aspects of the implementation more difficult.

## License

This code is released under the terms of the LGPLv3, but with an added exception that was shamelessly taken from [godeb](https://github.com/niemeyer/godeb).
Under certain circumstances, this exception permits distribution of binaries that are (statically or dynamically) linked with this code, without requiring the distribution of Minimal Corresponding Source or Minimal Application Code.
For more details, see: [LICENSE](LICENSE).
