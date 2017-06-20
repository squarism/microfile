Main features
==============
- YML configuration file mapping actions to URLs.  For example:

  ---
  default_url: http://localhost:3000
  triggers:
    /var/www/uploads/resumes:
      - /api/pdf_convert



Components
==========
- Configuration Manager
  - loads configuration from YAML file, supports one-to-many triggers
- Watcher
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
1. Filesystem listener (https://github.com/howeyc/fsnotify)
  - import "gopkg.in/fsnotify.v1"  (new API, recommended to use this as
    it's going into the stdlib)

1. Net/http lib (http://www.gorillatoolkit.org/pkg/http).  Or not if stdlib is fine.

1. Some way to run it.
desc "watch file for changes"
task :watch do
  require './lib/dropboy.rb'
  t = Dropboy::Trigger.new
  t.start
end

task :default => :watch

1. A way to test the filesystem watch?  How?
2. Generic handler for things other than web?  The config file would have to have like handler_type or something in it.
3. Recursive is a problem.  If supported, protection is needed.  What if "/" is watched?
4. Certain events are ignored, like chmod.
5. Sending the file contents is probably useful.  Need a switch for that.
6. If you do go recursive and you have a rules definition based on a filesystem hierarchy then you'll need to
   have a hierarchical rule system dude.  :(  Would child rules override children?
