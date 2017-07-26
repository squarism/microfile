# Dropboy

> It's your dropbox folder's best friend.

Dropboy watches a dropbox or folder for events and does something more meaningful with that event.  It's a
daemon designed to be everything you need to put next to an upload folder.

_This project is not affiliated or directly releated to Dropbox Inc. or their service.  It can be used in conjunction with a Dropbox folder._


## Neato Features

* Super readable Hashicorp configuration file.
* Available as a binary.  Easy to deploy!  <3 ops
* Action ordering is guaranteed, you could possibly create a tiny workflow of steps to custom APIs.
* Evented reactions, no polling the disk waiting for changes.
* Sane and structured logging.
* Gets your problem away from possibly a cheap fileserver and into a higher level of abstraction.


## Up and Running
1. Download a binary release.
2. Create a config file (more realistic examples are below) in one of three places:
    1. `$HOME/.dropboy/dropboy.hcl`
    2. `/usr/local/etc/dropboy.hcl`
    3. `./dropboy.hcl` (next to the binary `dropboy`) you downloaded.
3. Start Dropboy: `./dropboy`

See below for what actually goes in the config file!


## Use Cases and Examples

### Image Shrinking
Let's say you have a dropbox for clients to put pictures in.  And they upload huge images off their camera.  And you'd rather do something else than write a script and cron it up to shrink down the pictures.

Create a basic _Dropboy_ config file in `~/.dropboy/dropboy.hcl`.
```
log_file = "/var/log/dropboy.log"
default_url = "http://localhost:9000"

watch "/var/www/uploads/real_estate_photos" {
  action "imaginary" {
    options {
      output_directory = "/var/www/uploads/converted_images"
      path = "/resize?width=128&type=jpeg&nocrop=true"
    }
  }
}
```

[Imaginary](https://github.com/h2non/imaginary) is a microservice that can handle image processing with very few dependencies.  You can run imaginary very easily with Docker: `docker run -p 9000:9000 h2non/imaginary`.

Images will show up shrunk down and converted to `jpeg` in `/var/www/uploads/converted_images`.  Limitations of Imaginary are still present though.  It cannot handle all image types currently.

### New File Alert
Let's say you have a folder where resumes show up and you'd like to notify Slack or something.
If you have an API that can accept a `POST` and then do the Slack notification part then Dropboy is the perfect glue service.

```
default_url = "http://localhost:9000"

watch "/var/www/uploads/real_estate_photos" {
  action "imaginary" {
    options {
      output_directory = "/var/www/uploads/converted_images"
      path = "/resize?width=128&type=jpeg&nocrop=true"
    }
  }
}
```


## Options
Options go at the top level of the config file.

| Option | Use | Example |
|---|---|---|
| `default_url` | Reduces the amount of typing.  Partial paths in actions will use this value. |  "http://localhost:9000" |
| `log_file` | File to write logs to.  If left blank, goes to `stdout`. |
| `watch` | A single folder to watch.  Recursion not implemented right now.  You can watch multiple things. |


```
default_url = "http://resume-alerting.service:9000/api"
log_file = "/var/log/dropboy.log"

watch "/tmp/dropbox/resumes" {
  action "http" {
    options {
      path = "/new_resume"
    }
  }
}
```
This example assumes that the "resume-alerting service" only is going to alert based on the filename that just
showed up.


## Actions
You can have multiple actions per watch.  In that case, the actions happen in order of the config file
although it shouldn't matter too much since Dropboy doesn't have complex workflow features.

### Log
Simply print a message when a file changes
Syntax: `action "log"`

##### Additional Options
None


### Http
Send a custom `POST` to your service
Config Option: `action "http"`

##### Additional Options
| Option | Use |
|---|---|
| path | Full or partial URL path to send a message to |
| send_contents | Send the file contents as base64?  ("true" or "false") |


### Imaginary
Send the changed file to the imaginary microservice to do image operations
Config Option: `action "imaginary"`


##### Additional Options
| Option | Use | Example |
|---|---|
| path| URL of imaginary [operation](https://github.com/h2non/imaginary#http-api) endpoint.  Can be full or
partial url | `/resize | http://localhost:9000/resize?width=256` |
| output_directory | Where to store processed images | `"/var/www/processed_images"` |


## A Custom API Example
Ok, so what if you want to write an API?  Here's a rough example of how you could process a file coming from
Dropboy in Rails.

Create a controller (in this case `thumbnail_controller` and handle a `POST`.

```ruby
def create
  if !params[:thumbnail]
    head :no_content
    return
  end

  uploaded_filename = write_tmp_file(params[:thumbnail][:filename], params[:thumbnail][:contents])
  png_destination_file = File.basename(uploaded_filename.gsub(/\.jpg$/, '.png')).downcase
  destination_file = "#{Rails.root}/tmp/thumbnails/#{png_destination_file}"

  convert_to_thumbnail(uploaded_filename, destination_file)

  head :ok
end

private
def write_tmp_file(tmp_filename, contents)
  uploaded_filename = nil
  if contents
    basename = File.basename(tmp_filename)
    uploaded_filename = "#{Rails.root}/tmp/#{basename}"
    File.open(uploaded_filename, "wb") do |file|
      file.write Base64.decode64(contents)
    end
  end

  uploaded_filename
end

def convert_to_thumbnail(source, destination)
  convert_command = "convert \
    -thumbnail '256x256^' \
    -define png:size=256x256 \
    -gravity center \
    -extent 256x256 \
    -polaroid 3 \
    #{source} #{destination}"

  system(convert_command)
end
```
Dropboy will send contents to this API and then the thumbnail will show up in the Rails' `./tmp` folder.
Neat huh.  Note that this requires imagemagick to be installed because this API is just shelling out.
This is just a quick example, it's not an example of clean Rails code.


## Missing Features
* Recursion - hard to check and solve right now
* Catch up mode - file system changes that are missed when it's not running
* Wildcard filtering - if you want to just watch a path for `*.pdf`
* Chmod events are straight up ignored
* Templating - so you could watch for `*.exe` and have a Shell action `rm $file` or something.


## License
Dropboy is released under the [MIT License](http://www.opensource.org/licenses/MIT).


## Contributing
I encourage anyone to give feedback no matter what level.  If Dropboy is close to what you want but not quite,
create a pull request with no code!

Also, small suggestion for this project and any other: Please do not spend too much of you time on a fork or branch without
talking about it first.  It's much easier to pre-discuss an idea before a ton of code and time.
