cloudinary
===================

context base [go-cloudinary](https://github.com/gotsunami/go-cloudinary) wrapper for golang

[Cloudinary](http://cloudinary.com/)

## Usage

[example](https://github.com/kyokomi/cloudinary/blob/master/example/main.go)

```go
package main

import (
	"bytes"
	"io/ioutil"

	"github.com/kyokomi/cloudinary"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
    ctx = NewContext(ctx, "cloudinary://<API Key>:<API Secret>@<Cloud name>")

	data, _ := ioutil.ReadFile("<imageFile>")
	cloudinary.UploadStaticImage(ctx, "<name>", bytes.NewBuffer(data))
}
```

# License

[MIT](https://github.com/kyokomi/cloudinary/blob/master/LICENSE)
