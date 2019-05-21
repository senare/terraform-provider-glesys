resource "glesys_loadbalancer" "mylb" {
  count = 1
  datacenter = "Falkenberg"
  name = "mylb-1"
}
