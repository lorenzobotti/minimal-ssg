[date: 2020-08-02]
# Second post

![Image](https://justyy.com/wp-content/uploads/2016/01/markdown-syntax-language.png)

I'm working on the home page, and making a list of posts in chronological order for the blog. I'm also thinking of putting five or six of the most recet posts at the bottom of every post. I'm just worried about how to encode the pubblication date inside of the markdown file. Maybe make the first line into a timestamp?

## Todos

I still have to finish the markdown parser by adding list support and multiline code. Regular expressions are not enough for those because they would need to work on multiple lines at a time. Lists especially are going to be a pain because of indentation. I'll figure something out, but it's probably not going to take too much time.
I should also probably make the program move locally stored images from the `posts/` to the `output/` folder on its own.
The important part here is the Home Page. Once that's done, I'll basically have a respectable minimum viable product on my hands, and an almost finished one at that: I don't care for all the extra features, this is just an exercise, there's plenty of much better and more polished SSGs out there.