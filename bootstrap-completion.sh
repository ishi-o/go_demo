_bootstrap_complete() {
    local cur prev words cword
    _init_completion || return

    if [[ ${cword} -eq 1 ]]; then
        local modules=()
        while IFS= read -r line; do
            if [[ $line =~ ^[[:space:]]*\./([^[:space:]]+) ]]; then
                modules+=("${BASH_REMATCH[1]}")
            fi
        done < <(grep "^\s*\./" go.work 2>/dev/null)

        COMPREPLY=($(compgen -W "${modules[*]}" -- "${cur}"))
    fi
}

complete -F _bootstrap_complete bootstrap.sh

