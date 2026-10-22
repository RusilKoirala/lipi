import { useEffect, useState } from "react";
import { Copy, Download } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { GITHUB_URL, INSTALL_CMD, WINDOWS_SETUP_URL } from "@/lib/site";

function Command({ cmd }: { cmd: string }) {
  return (
    <div className="flex min-w-0 items-center gap-3 rounded-lg border bg-paper py-1.5 pr-1.5 pl-3.5 font-mono text-[13px]">
      <span className="text-crimson select-none" aria-hidden>
        $
      </span>
      <code className="min-w-0 flex-1 overflow-x-auto whitespace-nowrap">{cmd}</code>
      <button
        type="button"
        data-copy={cmd}
        aria-label={`Copy: ${cmd}`}
        className="inline-flex h-7 shrink-0 items-center gap-1.5 rounded-md px-2 font-sans text-xs text-muted-foreground transition-colors hover:bg-muted hover:text-foreground data-[copied]:text-emerald-600 dark:data-[copied]:text-emerald-400"
      >
        <Copy className="size-3.5" />
        <span data-copy-label>Copy</span>
      </button>
    </div>
  );
}

function Step({ n, children }: { n: string; children: React.ReactNode }) {
  return (
    <li className="grid grid-cols-[1.75rem_minmax(0,1fr)] gap-x-2">
      <span className="font-deva-text text-lg leading-6 font-medium text-crimson" aria-hidden>
        {n}
      </span>
      <div className="grid min-w-0 gap-2.5 text-[15px] leading-6">{children}</div>
    </li>
  );
}

export default function InstallTabs() {
  const [tab, setTab] = useState("unix");

  // Open the right tab for the visitor's computer.
  useEffect(() => {
    const ua = navigator.userAgent;
    if (/Windows/i.test(ua)) setTab("windows");
  }, []);

  return (
    <Tabs value={tab} onValueChange={setTab} className="mt-6">
      <TabsList>
        <TabsTrigger value="unix" className="px-3">
          macOS &amp; Linux
        </TabsTrigger>
        <TabsTrigger value="windows" className="px-3">
          Windows
        </TabsTrigger>
      </TabsList>

      <TabsContent value="unix" className="mt-6">
        <ol className="grid gap-7">
          <Step n="१">
            <p>Open Terminal, paste this and press Enter:</p>
            <Command cmd={INSTALL_CMD} />
            <p className="text-sm text-muted-foreground">
              It downloads the right lipi for your computer (Intel or Apple Silicon, x86 or ARM), checks the download,
              puts it in <code className="font-mono text-foreground">~/.lipi/bin</code> and adds that folder to your
              PATH. No password needed. Run the same command again later to update.
            </p>
          </Step>
          <Step n="२">
            <p>Open a new terminal window, then check it worked:</p>
            <Command cmd="lipi --version" />
          </Step>
        </ol>
        <p className="mt-6 text-sm text-muted-foreground">
          Want to read the script before running it?{" "}
          <a
            href={`${GITHUB_URL}/blob/master/install.sh`}
            className="text-foreground underline decoration-crimson underline-offset-2"
          >
            It is on GitHub
          </a>
          .
        </p>
      </TabsContent>

      <TabsContent value="windows" className="mt-6">
        <ol className="grid gap-7">
          <Step n="१">
            <p>Download the setup wizard:</p>
            <a
              href={WINDOWS_SETUP_URL}
              className="inline-flex h-10 w-fit items-center gap-2 rounded-lg bg-primary px-4 text-sm font-medium text-primary-foreground transition-opacity hover:opacity-90"
            >
              <Download className="size-4" />
              Download lipi-setup.exe
            </a>
            <p className="text-sm text-muted-foreground">Windows 10 or 11, 64-bit.</p>
          </Step>
          <Step n="२">
            <p>
              Open it and click through the wizard. Leave{" "}
              <span className="font-medium">“Add lipi to PATH”</span> ticked. It installs for your account only, so it
              doesn't ask for an administrator password.
            </p>
            <p className="rounded-lg border bg-paper px-3.5 py-2.5 text-sm text-muted-foreground">
              If Windows shows <span className="text-foreground">“Windows protected your PC”</span>, click{" "}
              <span className="text-foreground">More info</span>, then{" "}
              <span className="text-foreground">Run anyway</span>. This appears for new programs that aren't signed with
              a paid certificate yet.
            </p>
          </Step>
          <Step n="३">
            <p>
              Open a <span className="font-medium">new</span> Command Prompt or PowerShell window and check it worked:
            </p>
            <Command cmd="lipi --version" />
          </Step>
        </ol>
        <p className="mt-6 text-sm text-muted-foreground">
          To remove Lipi, go to Settings → Apps → Installed apps → Lipi → Uninstall. That also takes it off your PATH.
        </p>
      </TabsContent>
    </Tabs>
  );
}
