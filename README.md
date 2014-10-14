Main features
==============
- YML configuration file mapping actions to URLs.  For example:

  ---
  default_url: http://localhost:3000
  triggers:
    - ./spec/tmp/dropbox:
    - /api/run_loader



Components
==========
- Configuration Manager
  - loads configuration from YAML file, supports one-to-many triggers
- Trigger Object
  - holds list of watches
  - register_watch()
- Notifier
  - notify_webservice(url, message)
  - start(); stop()
- App
  - load configuration
  - main loop (run notifier)


Need
====
1. YAML reader

2. Filesystem listener (https://github.com/howeyc/fsnotify)
  - import "gopkg.in/fsnotify.v1"  (new API, recommended to use this as
    it's going into the stdlib)

3. Net/http lib (http://www.gorillatoolkit.org/pkg/http)

4. Some way to run it.
desc "watch file for changes"
task :watch do
  require './lib/dropboy.rb'
  t = Dropboy::Trigger.new
  t.start
end

task :default => :watch

5. Dependencies with Godeps (bundler)
