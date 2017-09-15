![image](https://github.com/squarism/microfile/raw/images/images/Microfile.png)

Microfile relays file events and/or file contents to microservices.  It is easily configurable
with simple rules and is easily deployable.  It's power comes from how you wire it up.


## Features

* Super readable Hashicorp configuration file.
* Available as a binary.  Easy to deploy!  <3 ops
* Action ordering is guaranteed, you could possibly create a tiny workflow of steps to custom APIs.
* Evented reactions, no polling the disk waiting for changes.
* Sane and structured logging.
* Gets your problem away from possibly a cheap fileserver and into a higher level of abstraction.


## Up and Running
1. Download a binary release.
2. Create a config file (more realistic examples are below) in one of three places:
    1. `$HOME/.microfile/microfile.hcl`
    2. `/usr/local/etc/microfile.hcl`
    3. `./microfile.hcl` (next to the binary `microfile`) you downloaded.
3. Start Microfile: `./microfile`

See below for what actually goes in the config file.


## Use Cases and Examples

### Image Shrinking

![image](https://github.com/squarism/microfile/raw/images/images/microfile_image_api.png)

Let's say you have a dropbox for clients to put pictures in.  And they upload huge images off their camera.  And you'd rather do something else than write a script and cron it up to shrink down the pictures.

Create a basic _Microfile_ config file in `~/.microfile/microfile.hcl`.
```
log_file = "/var/log/microfile.log"
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
If you have an API that can accept a `POST` and then do the Slack notification part then Microfile is the perfect glue service.

```
watch "/var/www/uploads/real_estate_photos" {
  action "http" {
    options {
      path = "http://your-api.local/file_arrival"
    }
  }
}
```


### Shared Folder Cleanup

Let's say you have a Dropbox or a Google Drive folder shared among multiple things.
Images arrive in png format but that wastes space.  They also aren't named with timestamps.
You could easily write an API to rename to the current time.  Microfile can convert
the image, drop it in a Microfile watched folder that can send an event to an API
that renames the file.

This is just an example of a simple workflow.

```
watch "/Users/you/Google Drive/screenshots" {
  action "imaginary" {
    options {
      output_directory = "/Users/you/Google Drive/screenshots/jpegs"
      path = "/convert?type=jpeg&quality=8"
    }
  }
}

watch "/Users/you/Google Drive/screenshots/jpegs" {
  action "http" {
    options {
      path = "http://your-api.local/timestamp_rename"
    }
  }
}
```

Because the folder is shared the server at `your-api.local` can take the file event and rename it.
Microfile _could_ do the rename but that's not the point of composition.  In the future there may
be other adapters and actions that Microfile takes on but it's mostly going to be around
message passing and not work doing.



## Options
Options go at the top level of the config file.

| Option | Use | Example |
|---|---|---|
| default_url | Reduces the amount of typing.  Partial paths in actions will use this value. |  `http://localhost:9000` |
| `log_file` | File to write logs to.  If left blank, goes to `stdout`. | `/tmp/microfile.log` |
| `watch` | A single folder to watch.  You can watch multiple things.  No need to escape paths.  Shell expansion like `~` does not work.  No recursion support right now.  | `watch "/var/foo"` |


```
default_url = "http://resume-alerting.service:9000/api"
log_file = "/var/log/microfile.log"

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
although it shouldn't matter too much since Microfile doesn't have complex workflow features.

### Log
Simply print a message when a file changes
Syntax: `action "log"`

##### Additional Options
None


### Http
Send a custom `POST` to your service
Config Option: `action "http"`

##### Additional Options
| Option | Use | Example |
|---|---|---|
| path | Full or partial URL path to send a message to |
| send_contents | Send the file contents as base64?  ("true" or "false") |


### Imaginary

Send the changed file to the imaginary microservice to do image operations
Config Option: `action "imaginary"`


##### Additional Options
| Option | Use | Example |
|---|---|---|
| path| URL of imaginary [operation](https://github.com/h2non/imaginary#http-api) endpoint.  Can be full or
partial url | `/resize | http://localhost:9000/resize?width=256` | |
| output_directory | Where to store processed images | `"/var/www/processed_images"` | |


## A Custom API Example

Ok, so what if you want to write an API?  Here's a rough example of how you could process a file coming from
Microfile in Rails.

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
Microfile will send contents to this API and then the thumbnail will show up in the Rails' `./tmp` folder.
Neat huh.  Note that this requires imagemagick to be installed because this API is just shelling out.
This is just a quick example, it's not an example of clean Rails code.


## Missing Features and TODO

These are coming ...

* Catch up mode - file system changes that are missed when it's not running
* Wildcard filtering - if you want to just watch a path for `*.pdf`
* Chmod events are straight up ignored
* Templating - so you could watch for `*.exe` and have a Shell action `rm $file` or something.  Or post to an existing API (json template).
* Recursive watching of folders - hard to check and solve right now


## Philosophy

Microfile is about passing message and not doing work itself.  This gets the problem
away from possibly weak machines and into a full language powered API or service.  Microfile is
defined by a config file which will never be as powerful and flexible as a language.

For example:
- image conversion is done by passing what's needed to a service.
- file renaming is done by passing the minimum of what's need to rename a file.
- testing a backup archive is done somewhere else even though it's easy to imagine a shell-out.


## License

Microfile is released under the [MIT License](http://www.opensource.org/licenses/MIT).


## Contributing
I encourage anyone to give feedback no matter what level.  If Microfile is close to what you want but not quite,
create a pull request with no code!

Also, small suggestion for this project and any other: Please do not spend too much of you time on a fork or branch without
talking about it first.  It's much easier to pre-discuss an idea before a ton of code and time.
