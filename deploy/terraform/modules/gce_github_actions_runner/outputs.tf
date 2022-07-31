output "ip" {
  description = "The IP address of the runner's VM."
  value       = google_compute_address.runner.address
}
