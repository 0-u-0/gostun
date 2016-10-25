Usage
------
```
usage: gostun [<flags>]

实时猫 TURN/STUN 服务器

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=config.json,/etc/rtcat_golang_stun/config.json
                                 Configuration file location
  -x, --external_ip=EXTERNAL_IP  TURN Server public/private address mapping, if the server is behind NAT.
      --min_port=49152           Lower bound of the UDP port range for relay endpoints allocation.
      --max_port=65535           Upper bound of the UDP port range for relay endpoints allocation.
  -v, --version                  Show application version.
```


RFC
-----

- [x] RFC 5389 : Session Traversal Utilities for NAT (STUN)
- [ ] RFC 5769 : Test Vectors for Session Traversal Utilities for NAT (STUN)
- [ ] RFC 5766 : Traversal Using Relays around NAT (TURN)
- [ ] RFC 5245 : Interactive Connectivity Establishment (ICE)

Features
--------

- [x] STUN UDP
- [x] Check Nonce
- [x] NAT Mapped Support
- [ ] Port Range
- [ ] FINGERPRINT
- [ ] STUN TCP
- [ ] STUN TCP over TLS
- [ ] TURN TCP
- [ ] TURN TCP over TLS
- [ ] IPV6


Test
------

- [ ] STUN Public Test


Expire
---------

- Nonce : 20 min
- Allocate : min : 5 min , max : 10 min