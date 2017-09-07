# Flappy Bird

This is my Flappy Bird attempt, which I created in the span of a half weekend, so there might be some rough edges.

The original plan is to build a genetic algorithm on top of this 🤖

But until then, we can enjoy the game.

## Preview

![Alt Text](https://github.com/Opposition/flappy-bird/raw/master/preview/flappy-bird.gif)

## Dependencies

### Compiling

To compile the game, you will need Git, Go and SDL2 installed. 

Follow the instructions for Git [here](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git). 
Follow the instructions for Go [here](https://golang.org/dl/).
Follow the instructions for SDL2 [here](https://github.com/veandco/go-sdl2).

After that has been done, run the following command:

    go get github.com/Opposition/flappy-bird

The generated binary should be in `$GOPATH/bin`.

### Playing

To play the game, the client will need SDL2 installed on their machine. Follow the instructions [here](https://github.com/veandco/go-sdl2).

Note that the data folder needs to follow the binary:

    flappy-bird/
    ├── data/
    │   ├── background/
    │       ├── background.png
    │   ├── bird/
    │       ├── bird_frame_1.png
    │       ├── bird_frame_2.png
    │       ├── bird_frame_3.png
    │       ├── bird_frame_4.png
    │       ├── bird_frame_5.png
    │       ├── bird_frame_6.png
    │       ├── bird_frame_7.png
    │       ├── bird_frame_8.png
    │       ├── bird_frame_9.png
    │       ├── bird_frame_10.png
    │   ├── font/
    │       ├── font.ttf
    │   ├── pipe/
    │       ├── pipe.png
    │   ├── sound/
    │       ├── intro.ogg
    └── .flappy-bird

## Assets

All the images, fonts and songs used in this game are Creative Commons licensed. They are obtained from [here](https://opengameart.org/).

This project is licensed under Apache v2 license. Read more in the [LICENSE](LICENSE) file.