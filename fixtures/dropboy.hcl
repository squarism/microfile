default_url = "http://localhost:9876"

watch "/var/www/resumes/dropbox" {
  actions "http" {
    options {
      send_file = "true"
    }
  }
}
