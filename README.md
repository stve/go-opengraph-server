# go-opengraph-server

A small microservice to retrieve Open Graph data from websites and return as JSON.

# Usage

Pass the URL you want to extract data from by passing the `url` parameter:

```
http://localhost:8000/?url=https://vimeo.com/205090959
```

Parsed data is returned as JSON:

```json
{
  "site_name": "Vimeo",
  "type": "video",
  "title": "Symphony of Light - Kauai Timelapse",
  "description": "From the towering green spires to the pristine beaches, the stunning island of Kauai offers an incredible range of unique landscapes to explore. \"Symphony of…",
  "url": "https://vimeo.com/205090959",
  "videos": [
    {
      "height": "545",
      "secure_url": "https://player.vimeo.com/video/205090959?autoplay=1",
      "type": "text/html",
      "url": "https://player.vimeo.com/video/205090959?autoplay=1",
      "width": "1280"
    },
    {
      "height": "545",
      "secure_url": "https://vimeo.com/moogaloop.swf?clip_id=205090959\u0026autoplay=1",
      "type": "application/x-shockwave-flash",
      "url": "https://vimeo.com/moogaloop.swf?clip_id=205090959\u0026autoplay=1",
      "width": "1280"
    }
  ],
  "images": [
    {
      "height": "545",
      "secure_url": "https://i.vimeocdn.com/video/619779422_1280x545.jpg",
      "type": "image/jpg",
      "url": "https://i.vimeocdn.com/video/619779422_1280x545.jpg",
      "width": "1280"
    }
  ]
}
```

## Caveats

My use case was for fetching data from Vimeo and YouTube. I haven't tested other sites and some assumptions were made in generating JSON based on my specific needs that may not make sense or be applicable to other sites and their Open Graph data.

## Additional Caveats

I don't write much go code, feedback is **definitely** welcome!

## License

MIT (c) Steve Agalloco. See [LICENSE](https://github.com/stve/go-opengraph-server/blob/master/LICENSE.md) for details.
