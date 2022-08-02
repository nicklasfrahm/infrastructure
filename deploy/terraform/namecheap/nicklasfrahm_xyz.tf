resource "namecheap_domain_records" "nicklasfrahm_xyz" {
  domain = "nicklasfrahm.xyz"

  # dktil01
  record {
    type     = "CNAME"
    hostname = "dktil01"
    address  = "alfa.nicklasfrahm.xyz"
  }
  record {
    type     = "CNAME"
    hostname = "*.dktil01"
    address  = "alfa.nicklasfrahm.xyz"
  }

  # deflf01
  record {
    type     = "CNAME"
    hostname = "deflf01"
    address  = "bravo.nicklasfrahm.xyz"
  }
  record {
    type     = "CNAME"
    hostname = "*.deflf01"
    address  = "bravo.nicklasfrahm.xyz"
  }

  # deflf02
  record {
    type     = "CNAME"
    hostname = "deflf02"
    address  = "charlie.nicklasfrahm.xyz"
  }
  record {
    type     = "CNAME"
    hostname = "*.deflf02"
    address  = "charlie.nicklasfrahm.xyz"
  }

  # uscbf01
  record {
    type     = "CNAME"
    hostname = "uscbf01"
    address  = "delta.nicklasfrahm.xyz"
  }
  record {
    type     = "CNAME"
    hostname = "*.uscbf01"
    address  = "delta.nicklasfrahm.xyz"
  }
}
