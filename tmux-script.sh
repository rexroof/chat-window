#!/bin/bash

# our go code makes a pane with this title
export PANE_TITLE="chatwindowaeP8e"
HEIGHT=6

# find an existing pane where our go code is running
PANE_ID=$(tmux list-panes -as -f "#{==:#{pane_title},$PANE_TITLE}" -F '#S:#I.#D')

# if we were passed the arg, run our go code now.
if [ "${1}" == "intmux" ] ; then
  cd $HOME/go/src/github.com/rexroof/chat-window
  source .env
  # go run main.go
  /home/rex/bin/chatwindow
elif [ "${1}" == "send" ] ; then
  # shift removes the "send" from our cmd line
  shift

  if [ -n "${PANE_ID}" ] ; then
    # send keys to the chatwindow pane, with newline
    tmux send-keys -t "${PANE_ID}" "$@" C-m
  else
    echo "cannot find chatwindow pane"
    exit 1
  fi

else
  if [ -n "${PANE_ID}" ] ; then
    # pull that pane to our current window
    tmux join-pane -vfdl $HEIGHT -s "${PANE_ID}"
  else
    # or start a new pane with our go code!
    tmux split-window -vfdl $HEIGHT /home/rex/bin/chat-window.sh intmux
  fi

  exit 0
fi
