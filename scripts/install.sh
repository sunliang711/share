#!/bin/bash
rpath="$(readlink ${BASH_SOURCE})"
if [ -z "$rpath" ];then
    rpath=${BASH_SOURCE}
fi
pwd=${PWD}
this="$(cd $(dirname $rpath) && pwd)"
# cd "$this"
export PATH=$PATH:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

user="${SUDO_USER:-$(whoami)}"
home="$(eval echo ~$user)"

# export TERM=xterm-256color

# Use colors, but only if connected to a terminal, and that terminal
# supports them.
if which tput >/dev/null 2>&1; then
  ncolors=$(tput colors 2>/dev/null)
fi
if [ -t 1 ] && [ -n "$ncolors" ] && [ "$ncolors" -ge 8 ]; then
    RED="$(tput setaf 1)"
    GREEN="$(tput setaf 2)"
    YELLOW="$(tput setaf 3)"
    BLUE="$(tput setaf 4)"
            CYAN="$(tput setaf 5)"
    BOLD="$(tput bold)"
    NORMAL="$(tput sgr0)"
else
    RED=""
    GREEN=""
    YELLOW=""
            CYAN=""
    BLUE=""
    BOLD=""
    NORMAL=""
fi
_err(){
    echo "$*" >&2
}

_runAsRoot(){
    cmd="${*}"
    local rootID=0
    if [ "${EUID}" -ne "${rootID}" ];then
        echo -n "Not root, try to run as root.."
        # or sudo sh -c ${cmd} ?
        if eval "sudo ${cmd}";then
            echo "ok"
            return 0
        else
            echo "failed"
            return 1
        fi
    else
        # or sh -c ${cmd} ?
        eval "${cmd}"
    fi
}

rootID=0
function _root(){
    if [ ${EUID} -ne ${rootID} ];then
        echo "Need run as root!"
        exit 1
    fi
}

ed=vi
if command -v vim >/dev/null 2>&1;then
    ed=vim
fi
if command -v nvim >/dev/null 2>&1;then
    ed=nvim
fi
if [ -n "${editor}" ];then
    ed=${editor}
fi
###############################################################################
# write your code below (just define function[s])
# function is hidden when begin with '_'
###############################################################################
# TODO

serviceName="share-srv"
cliName="share-cli"
installServer(){
    cd "${this}"
    _runAsRoot "systemctl stop ${serviceName}.service"
    build
    local dest="${1:?'missing install location'}"
    if [ ! -e config.toml ];then
        echo "Not found config.toml" >&2
        exit 1
    fi
    if [ ! -e ${serviceName} ];then
        echo "Not found ${serviceName},build it first" >&2
        exit 1
    fi

    if [ ! -d "${dest}" ];then
        echo -n "Create dir: ${dest} ..."
        mkdir -p "${dest}" || { echo "failed!"; exit 1; }
    fi

    mkdir -p "share-server" || { echo "Create share-server failed!"; exit 1; }

    cp config.toml "${dest}"
    cp ${serviceName} "${dest}"

    ${ed} ${dest}/config.toml

    sed -e "s|DEST|$dest|g" \
        -e "s|EXE|${serviceName}|g" \
        ${serviceName}.service > /tmp/${serviceName}.service
    _runAsRoot "mv /tmp/${serviceName}.service /etc/systemd/system/${serviceName}.service"
    _runAsRoot "systemctl enable --now ${serviceName}.service"
}

build(){
    flags="-X main.buildstamp=`date +%FT%T` -X main.githash=`git rev-parse HEAD`"
    echo "build ${serviceName}..."
    cd "${this}"
    go build -ldflags "$flags" -o ${serviceName} "${this}/../cmd/share"

    echo "build ${cliName}..."
    go build -ldflags "$flags" -o ${cliName} "${this}/../cmd/cli"
}

installClient(){
    local dest="${1:-/usr/local/bin}"
    cd "${this}"
    cp ${cliName} "${dest}"
    cp .share.toml $home
    ${ed} $home/.share.toml
}

em(){
    $ed $0
}

###############################################################################
# write your code above
###############################################################################
function _help(){
    cat<<EOF2
Usage: $(basename $0) ${bold}CMD${reset}

${bold}CMD${reset}:
EOF2
    # perl -lne 'print "\t$1" if /^\s*(\w+)\(\)\{$/' $(basename ${BASH_SOURCE})
    # perl -lne 'print "\t$2" if /^\s*(function)?\s*(\w+)\(\)\{$/' $(basename ${BASH_SOURCE}) | grep -v '^\t_'
    perl -lne 'print "\t$2" if /^\s*(function)?\s*(\w+)\(\)\{$/' $(basename ${BASH_SOURCE}) | perl -lne "print if /^\t[^_]/"
}

function _loadENV(){
    if [ -z "$INIT_HTTP_PROXY" ];then
        echo "INIT_HTTP_PROXY is empty"
        echo -n "Enter http proxy: (if you need) "
        read INIT_HTTP_PROXY
    fi
    if [ -n "$INIT_HTTP_PROXY" ];then
        echo "set http proxy to $INIT_HTTP_PROXY"
        export http_proxy=$INIT_HTTP_PROXY
        export https_proxy=$INIT_HTTP_PROXY
        export HTTP_PROXY=$INIT_HTTP_PROXY
        export HTTPS_PROXY=$INIT_HTTP_PROXY
        git config --global http.proxy $INIT_HTTP_PROXY
        git config --global https.proxy $INIT_HTTP_PROXY
    else
        echo "No use http proxy"
    fi
}

function _unloadENV(){
    if [ -n "$https_proxy" ];then
        unset http_proxy
        unset https_proxy
        unset HTTP_PROXY
        unset HTTPS_PROXY
        git config --global --unset-all http.proxy
        git config --global --unset-all https.proxy
    fi
}


case "$1" in
     ""|-h|--help|help)
        _help
        ;;
    *)
        "$@"
esac
