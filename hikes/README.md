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
* Publish: `hugo --environment production --minify`
* Local test: `python3 -m http.server 9000`
