# Hugo based hiking blog

## Credits

* https://gohugo.io/
* https://kaiiiz.github.io/hugo-theme-monochrome/

## Notes

* `hugo new site hikes`
* Copy the template to `themes`
* Update `hugo.toml`
* Add a blog post:
  ```
  hugo new content content/posts/columbia-river-gorge.md
  hugo new content content/posts/ecola-state-park.md
  ```
* Install npm dependencies via `npm install`
* Any change to `map.js`, need to re-generate minified via `npm run build:js`
* Add GPS coords, & do a `make run` to re-generate `static/map/locations.json` 
* Publish: `hugo --minify`
* Local test: `hugo server`

