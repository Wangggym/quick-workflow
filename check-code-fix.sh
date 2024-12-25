#!/bin/bash

check_code_fix() {
  if command -v make &> /dev/null; then
      if make -qp | grep -q 'fix:'; then
          echo "Start excution make fix"
          make fix
      else
          echo "Can not excution make fix"
      fi
  fi
}