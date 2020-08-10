[date: 2020-08-03]
# First minimum viable product!

The homepage is done. The homepage now shows for each post a full preview with the title linking to the actual post itself. This is fine for now, but will get heavy once we start serving more posts. I'm thinking I can just show just the first 5 or so posts as full previews and the rest as just the title with a link to the post, maybe in their own "archive" page, or maybe just in the homepage below the previews.

## What now?
Well, I still need to add full markdown support (lists and code blocks) and maybe polish the style a bit, but I think we have a (**very** minimalistic) fully functioning static site generator on our hands. Actually, let's write what I need to do in the form of a list, so when I'm done implementing those this part will actually make sense:
 * List and code support for markdown
 * Locally stored images support (automatically move them form `posts/` to `output/` without any collisions)
 * Name collision detection for posts(?)
 * Some way to include the date of pubblication in the 
 * Probably something else that I can't remember right now, I'll add more stuff fo sure

## How am i gonna use this?
I think I'm just going to use this as a general mind journal, maybe detailing my personal projects (whether related to computer science or not), but always as if I was writing on an actual blog. Who knows? Maybe I'll actually pubblish this someday. Of course, if I manage to keep writing at an acceptable rate, which we all know I will not.