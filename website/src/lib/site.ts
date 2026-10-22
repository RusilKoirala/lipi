export const GITHUB_URL = "https://github.com/RusilKoirala/lipi";
export const GITHUB_REPO_API = "https://api.github.com/repos/RusilKoirala/lipi";
export const BOOK_URL = "https://interpreterbook.com";

/** macOS + Linux: one command that downloads lipi and adds it to PATH (see install.sh in the repo). */
export const INSTALL_CMD = "curl -fsSL https://raw.githubusercontent.com/RusilKoirala/lipi/master/install.sh | sh";
/** Windows: the setup wizard built by the release workflow. */
export const WINDOWS_SETUP_URL = `${GITHUB_URL}/releases/latest/download/lipi-setup.exe`;
/** For people who already have Go. */
export const GO_INSTALL_CMD = "go install github.com/rusilkoirala/lipi@latest";

/** Prefix a path with the site's base URL (needed if the site is served from a sub-folder). */
export function href(path: string): string {
  const base = import.meta.env.BASE_URL.replace(/\/$/, "");
  return `${base}${path.startsWith("/") ? path : `/${path}`}`;
}
