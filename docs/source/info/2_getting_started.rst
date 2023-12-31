.. _getting_started:

.. include:: ../_inc/head.rst

.. include:: ../_inc/in_progress.rst


##################
2- Getting Started
##################

Installation
############

Download the binary for your desired target system from the `latest relase page <https://github.com/superstes/calamary/releases/latest>`_

Run
###

Just execute it.

.. code-block:: bash

    /usr/bin/calamary


Config-validation only:

.. code-block:: bash

    /usr/bin/calamary -v


Modes
#####

See: :ref:`Modes <modes>`


Configuration
#############

See :ref:`Rules <rules>` for more details about defining the filter ruleset.

The default config path is :code:`/etc/calamary/config.yml`

Basic config example:

.. code-block:: yaml

    ---

    service:
      listen:
        - mode: 'transparent'
          ip4: ['127.0.0.1']  # is default
          ip6: ['::1']  # is default
          port: 4128
          tcp: true
          udp: false  # not yet implemented
          tproxy: false

      certs:
        serverPublic: '/tmp/calamary.crt'
        serverPrivate: '/tmp/calamary.key'
        interceptPublic: '/tmp/calamary.subca.crt'
        interceptPrivate: '/tmp/calamary.subca.key'

      dnsNameservers: ['1.1.1.1', '8.8.8.8']
      debug: false
      timeout:  # ms
        connect: 2000
        process: 1000
        idle: 30000

      output:
        retries: 1  # connect-retries

      metrics:
        enabled: false
        port: 9512

      vars:
        - name: 'net_private'
          value: ['192.168.0.0/16', '172.16.0.0/12', '10.0.0.0/8']
        - name: 'svc_http'
          value: [80, 443]

      rules:
        - match:
            dest: '192.168.100.0/24'
          action: 'drop'

        - match:
            port: ['!443', '!80']
          action: 'drop'

        - match:
            src: '$net_private'
            dest: '$net_private'
            port: '$svc_http'
            protoL4: 'tcp'
          action: 'accept'

        - match:
            dest: '!$net_private'
            port: 443
            protoL4: 'tcp'
          action: 'accept' 
 
 
Redirect traffic
################

See :ref:`Redirect traffic <redirect>`


Systemd service
===============

Example systemd service to run Calamary:


.. code-block:: text

    [Unit]
    Description=Calamary Forward Proxy
    Documentation=https://docs.calamary.net
    Documentation=https://github.com/superstes/calamary
    After=network-online.target
    Wants=network-online.target

    [Service]
    Type=simple
    Environment=CONFIG=/etc/calamary/config.yml

    # validate before start/restart
    ExecStartPre=/usr/bin/calamary -f $CONFIG -v
    ExecStart=/usr/bin/calamary -f $CONFIG

    # validate before reload
    ExecReload=/usr/bin/calamary -f $CONFIG -v
    ExecReload=/bin/kill -HUP $MAINPID

    User=proxy
    Group=proxy
    Restart=on-failure
    RestartSec=5s

    StandardOutput=journal
    StandardError=journal
    SyslogIdentifier=calamary

    [Install]
    WantedBy=multi-user.target


Build from sources
##################

Download and 'install' Golang 1.21 to build the binary from sources: `Golang download <https://go.dev/doc/install>`_

.. code-block:: bash

    git clone https://github.com/superstes/calamary
    cd calamary/lib
    go mod download
    cd main/
    go build -o calamary
