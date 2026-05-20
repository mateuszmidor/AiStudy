# OpenSpec

https://github.com/Fission-AI/OpenSpec/
- lightweight when compared to github speckit

## Install openspec tool

```sh
sudo pacman -S nvm
echo 'source /usr/share/nvm/init-nvm.sh' >> ~/.bashrc  # or ~/.zshrc
exec $SHELL
nvm install --lts
nvm use --lts
npm install -g @fission-ai/openspec@latest
```

## Work

```sh
cd yourproject/
openspec init # initialize it to use opencode
opencode
/opsx:propose <what-you-want-to-build> # create spec
/opsx:apply # implement all tasks from created spec
/opsx:archive # move to archive folder; development of this change is finished
#/opsx:explore - to explore ideas, ask questions, find alternatives
```

## General workflow

/opsx:propose ──► /opsx:apply ──► /opsx:sync ──► /opsx:archive