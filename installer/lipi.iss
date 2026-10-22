; Lipi setup wizard for Windows (Inno Setup 6).
;
; Built by .github/workflows/release.yml:
;   iscc /DAppVersion=0.1.0 installer\lipi.iss
; which expects out\lipi-windows-amd64.exe to exist and writes out\lipi-setup.exe.
;
; Installs for the current user only (no admin prompt) into
; %LOCALAPPDATA%\Programs\Lipi, adds its bin folder to the user PATH,
; and removes it again on uninstall.

#ifndef AppVersion
  #define AppVersion "0.0.0"
#endif
; Windows file versions must be plain numbers (0.1.0), so tags like v0.1.0-beta pass this separately.
#ifndef FileVersion
  #define FileVersion AppVersion
#endif

#define AppName "Lipi"
#define AppURL "https://github.com/RusilKoirala/lipi"

[Setup]
; Never change AppId: Windows uses it to recognise upgrades and the uninstaller.
AppId={{63614A76-6122-4DDC-85DC-E2CC981F2C0B}
AppName={#AppName}
AppVersion={#AppVersion}
AppVerName={#AppName} {#AppVersion}
AppPublisher=Rusil Koirala
AppPublisherURL={#AppURL}
AppSupportURL={#AppURL}/issues
AppUpdatesURL={#AppURL}/releases
VersionInfoVersion={#FileVersion}
VersionInfoDescription=Lipi programming language installer

PrivilegesRequired=lowest
DefaultDirName={autopf}\{#AppName}
DisableProgramGroupPage=yes
DisableReadyPage=no
UsePreviousAppDir=yes

ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible

ChangesEnvironment=yes
WizardStyle=modern
Compression=lzma2
SolidCompression=yes

OutputDir=..\out
OutputBaseFilename=lipi-setup

UninstallDisplayName={#AppName} {#AppVersion}
UninstallDisplayIcon={app}\bin\lipi.exe

[Messages]
WelcomeLabel2=This will install [name/ver] on your computer.%n%nLipi is a small programming language you write with Nepali words. After setup, open a terminal and type lipi to use it.
FinishedLabel=Lipi is installed.%n%nOpen a NEW Command Prompt or PowerShell window and type:%n%n    lipi --version%n%nAn example program is in the examples folder inside the install folder.

[Tasks]
Name: "addtopath"; Description: "Add lipi to PATH, so you can type lipi in any terminal (recommended)"

[Files]
Source: "..\out\lipi-windows-amd64.exe"; DestDir: "{app}\bin"; DestName: "lipi.exe"; Flags: ignoreversion
Source: "..\examples\*.lip"; DestDir: "{app}\examples"; Flags: ignoreversion skipifsourcedoesntexist

[Icons]
Name: "{autoprograms}\Lipi\Lipi website"; Filename: "{#AppURL}"
Name: "{autoprograms}\Lipi\Uninstall Lipi"; Filename: "{uninstallexe}"

[Registry]
; Append {app}\bin to the user's PATH (only if it is not there yet).
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}\bin"; Tasks: addtopath; Check: NeedsAddPath(ExpandConstant('{app}\bin'))

[Run]
Filename: "{#AppURL}"; Description: "Open the Lipi website"; Flags: postinstall shellexec nowait skipifsilent unchecked

[Code]
const
  EnvKey = 'Environment';

{ True when Dir is not already in the user's PATH. }
function NeedsAddPath(Dir: string): Boolean;
var
  Path: string;
begin
  if not RegQueryStringValue(HKEY_CURRENT_USER, EnvKey, 'Path', Path) then
  begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + Uppercase(Dir) + ';', ';' + Uppercase(Path) + ';') = 0;
end;

{ Remove Dir from the user's PATH. Called when Lipi is uninstalled. }
procedure RemovePath(Dir: string);
var
  Path: string;
  P: Integer;
begin
  if not RegQueryStringValue(HKEY_CURRENT_USER, EnvKey, 'Path', Path) then
    exit;
  P := Pos(';' + Uppercase(Dir) + ';', ';' + Uppercase(Path) + ';');
  if P = 0 then
    exit;
  { P points at the ';' before Dir in ';' + Path + ';', so Dir starts at Path[P]. }
  if P = 1 then
    Delete(Path, 1, Length(Dir) + 1)
  else
    Delete(Path, P - 1, Length(Dir) + 1);
  RegWriteExpandStringValue(HKEY_CURRENT_USER, EnvKey, 'Path', Path);
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
    RemovePath(ExpandConstant('{app}\bin'));
end;
