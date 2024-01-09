# MTG Deck Pic Generator
Skip finding your cards from a mountain of deckboxes/trade binders and **manually** aligning your cards before taking a perfect MTG deck picture. Start **automatically** generating MTG deck pictures with your decklists and your own card collection.

## Usage
Download this repo by clicking the following buttons:

![download](https://github.com/manted/mtg_deck_pic/blob/main/screenshots/download.png?raw=true)

Unzip it to a desired folder on your machine.

To run the program, please refer to the OS system specific to your machine:
- [Windows](#on-windows)
- [Mac](#on-mac)

#### On Windows
On Windows, search for terminal in Search and select Windows Terminal from the search results.

In the terminal, go to the unzipped folder by running:
```sh
cd /path/to/your/folder
```

For example:
```sh
cd C:\Users\manted\mtg_deck_pic
```

To run the program in the terminal:

```sh
start deck_pic_generator DECKNAME
```

For example:
```sh
start deck_pic_generator grixis_delver
```

#### On Mac
Open the Terminal app.

In the terminal, go to the unzipped folder by running:
```sh
cd /path/to/your/folder
```

For example:
```sh
cd /Users/manted/mtg_deck_pic
```

To run the program in the terminal:

```sh
./deck_pic_generator DECKNAME
```

For example:
```sh
./deck_pic_generator grixis_delver
```

The `DECKNAME` should be one of the text file names in the `/decklist` folder, without the `txt` extension. The generated image would be `/decklist/DECKNAME.jpg`.

### Card Names
The card names in the decklists must match exactly the image names in the `/img` folder, without the `png` extension.

The individual card images must be in `png` format.

### Deck Picture Background
The background of the generated image is `/img/background.jpg`. You can replace it with any image of your choice and regenerate the deck pictures.

## Limitations
- Main deck size has to be exactly 60 or 80. Sideboard size has to be exactly 15.
- Split card names with `/` in them are not working at the moment. eg: Fire // Ice. You can use a different name without the `/` in the decklists and image names in `/img` folder. eg: Fire Ice.
- Cards that are in the decklists but are not in the `/img` folder will be left blank in the generated deck picture.
