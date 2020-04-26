# Gnome3 Xinput Layout Switch

Switch Gnome3 keyboard layout on release of two keys combination  (Ctrl Shift by default).
It is aimed to workaround https://bugs.launchpad.net/ubuntu/+source/gnome-control-center/+bug/36812 

It uses *xinput* under the hood and runs gdbus call each time to switch between tho most recent layouts in gnome with respect to status bar indication
       
       gdbus call --session --dest org.gnome.Shell --object-path /org/gnome/Shell --method org.gnome.Shell.Eval "imports.ui.status.keyboard.getInputSourceManager()._mruSources[1].activate()" 

## Configuration
Has command line arguments
* *--debug* - dump all keyboard events to show key codes  
* *--key1* - first key code to monitor (default: 37,105 [Ctrl])
* *--key2* - second key code to monitor (default: 50,62 [Shift])

## Setup instructions:

The simple way to set it up with the default configuration by running following commands in terminal:
    
    $ git clone https://gitlab.com/softkot/gnome3_xinput_layout_switch.git
    
    $ cd gnome3_xinput_layout_switch
    
    $ go generate
    
Instead of build it from source yuu can download it from [release page](https://gitlab.com/softkot/gnome3_xinput_layout_switch/-/releases) and continue afterward.
    
    $ sudo cp gnome-xinput-layout-switch /usr/bin/gnome-xinput-layout-switch
    
    $ echo /usr/bin/gnome-xinput-layout-switch \& | sudo tee /etc/X11/Xsession.d/99-gnome-xinput-layout-switch 

Then remove or disable gnome builtin keyboard shortcuts and restart X11.   
 
P.S.

* In case you want to change layout switch to Alt + Shift pass *--key1 64,108* argument.

* In case you want to change layout switch to Ctrl + Alt pass *--key2 64,108* argument.