#!/usr/bin/env bash
set -e

INSTALL_DIR="${HOME}/.local/bin"
TARGET="${INSTALL_DIR}/dsu"

# -------- proxy (按需配置) --------
PROXY="${https_proxy:-http://127.0.0.1:7897}"
export https_proxy="${PROXY}"
export http_proxy="${PROXY}"

# -------- Go 环境检测 --------
if ! command -v go &>/dev/null; then
    GOROOT="${HOME}/go-toolchain/go"
    if [ -x "${GOROOT}/bin/go" ]; then
        export GOROOT PATH="${GOROOT}/bin:${PATH}"
    else
        echo "未找到 Go，正在下载..."
        GO_TAR="go1.21.13.linux-amd64.tar.gz"
        wget -q "https://go.dev/dl/${GO_TAR}" -O /tmp/${GO_TAR}
        mkdir -p "$(dirname "${GOROOT}")"
        rm -rf "${GOROOT}"
        tar -C "$(dirname "${GOROOT}")" -xzf /tmp/${GO_TAR}
        rm -f /tmp/${GO_TAR}
        export GOROOT PATH="${GOROOT}/bin:${PATH}"
    fi
fi

echo "Go 版本: $(go version)"

# -------- 构建 --------
cd "$(dirname "$0")"

echo "下载依赖..."
go mod tidy

echo "编译..."
go build -o dsu .

# -------- 安装到 PATH --------
mkdir -p "${INSTALL_DIR}"
cp dsu "${TARGET}"
chmod +x "${TARGET}"

# 添加到 PATH（幂等）
SHELL_RC=""
case "${SHELL}" in
    */zsh) SHELL_RC="${HOME}/.zshrc" ;;
    */bash) SHELL_RC="${HOME}/.bashrc" ;;
    *) SHELL_RC="${HOME}/.bashrc" ;;
esac

if ! grep -q "${INSTALL_DIR}" "${SHELL_RC}" 2>/dev/null; then
    echo "export PATH=\"${INSTALL_DIR}:\${PATH}\"" >> "${SHELL_RC}"
    echo "已添加 ${INSTALL_DIR} 到 PATH (${SHELL_RC})"
fi

echo ""
echo "安装完成。运行: dsu"
echo "新终端需: source ${SHELL_RC} 或重新打开终端"
echo ""
echo "首次使用设置 API Key:"
echo "  dsu apikey sk-xxxxx"
