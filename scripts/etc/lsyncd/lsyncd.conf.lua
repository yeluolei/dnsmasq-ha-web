settings {
  logfile = "/var/log/lsyncd.log",
  statusFile = "/var/log/lsyncd.status",
  statusInterval = 5,
  maxDelays = 1,
  inotifyMode = "CloseWrite or Modify",
  nodaemon = false,
  maxProcesses = 1,
}

sync {
  default.rsyncssh,
  source = "/etc/dnsmasq-ha-web",
  host = "ubuntu@10.10.11.11",
  targetdir = "/etc/dnsmasq",
  delay = 5,
  delete = "running",
  ssh = {
    identityFile = "/etc/lsyncd/id_rsa",
  },
  rsync = {
    compress = true,
    checksum = true,
    archive  = true,
    rsync_path = "sudo /usr/local/bin/dnsmasq-rsync.sh",
  }
}

--- if you have more than 1 remote host, add more sections like above
sync {
  default.rsyncssh,
  source = "/etc/dnsmasq-ha-web",
  host = "ubuntu@10.10.12.12",
  targetdir = "/etc/dnsmasq",
  delay = 5,
  delete = "running",
  ssh = {
    identityFile = "/etc/lsyncd/id_rsa",
  },
  rsync = {
    compress = true,
    checksum = true,
    archive  = true,
    rsync_path = "sudo /usr/local/bin/dnsmasq-rsync.sh",
  }
}
