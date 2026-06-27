# Spectra Assets

The Spectra Explorer requires some assets like images and validator data in some cases (Gnoland explicitly).
In the future this might hold more data as the explorer evolves.
All of these assets are free to use and gather for any of your purpose as long as you respect the validators,
their identity, and image as they are made available by the validators themselves. The same applies to the
any other assets or images you find here.

If you are Validator read down below if you want to apply any sort of special image.

## Validator Image and Details

Some chains, like Gnoland, do not have any sort of way to tie some personal data properly. For example, you 
can't use any sort of emoji in the moniker (like Cogwheel usually does with "Cogwheel ⚙️"). There is also the
problem on not being able to attach an image to the validator profile. And reading validator Markdown and 
parsing it is hard because each validator might insert data differently. For Cosmos SDK chains for now it is
allowed to insert an image into this assets repo. Weather you want to overwrite the image that is available
with the keybase identity, or you simply do not want to register to Keybase this is the place where you can
insert image.

Just a note that the data inserted here will only be present in the Spectra Explorer. In other apps and 
explorers this won't probably reflect unless they also use this assets repo.

### How to Add the Data

To add image and any data first fork this repo, if you do not know how to do this you can follow the instructions [here](https://docs.github.com/en/get-started/quickstart/fork-a-repo).

In the data directory find the chain where your validators operates and make a directory with a name that 
matches the validator address. Cosmos SDK chains use bench32valoper. For Gnoland use strictly use the 
operator address, do not use the signer address.

For images place your logo image in the directory. It must be named `logo.jpeg` or `logo.jpg`. The `.png` are
not yet supported. This might change in the future but for now use strictly use `.jpg`.

This only applies to Gnoland or Tendermint2 chains so skip this for Cosmos SDK chains.
For any additional data you can make validator.json file which can contain data like this:
```json
{
  "moniker": "Cogwheel ⚙️",
  "identity": "3AB058E3A1912759",
  "description": "Cogwheel is a PoS validator that offers to delegators explorer to manage your assets, restake-ui and restake services so you can compound your rewards! For devs and users who need it there are public endpoints available. More info on https://cogwheel.zone!",
  "website": "https://cogwheel.zone",
  "security_contact": "info@cogwheel.zone"
}
```

Moniker is the name that will be shown in the explorer. The identity refers to the Keybase ID of the 
validator. This is useful if you want to use the same image you use for other validators on the Gnoland. Do 
not use it if you plan to add logo.jpg into the directory. Description is the description that will be shown 
in the explorer. However limit yourself to 256 characters.
