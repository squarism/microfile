default_url = "http://localhost:9876"

watch "/var/www/resumes/dropbox" {
  action "http" {
    options {
      send_file = "true"
    }
  }
}
