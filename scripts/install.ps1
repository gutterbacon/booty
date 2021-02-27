$repo = "amplify-edge/booty"
$filenamePattern = "*windows_amd64.zip"

$pathExtract = "C:\booty"

$releaseUrl = "https://api.github.com/repos/$repo/releases/latest"

function download_latest()
{
    Write-Host "getting latest release"
    $downloadUrl = ((Invoke-RestMethod -Method GET -Uri $releaseUrl).assets | Where-Object name -like $filenamePattern).browser_download_url
    $pathZip = Join-Path -Path $([System.IO.Path]::GetTempPath() ) -ChildPath $( Split-Path -Path $downloadUrl -Leaf )
    Write-Host "download latest release: $downloadUri to $pathZip"
    Invoke-WebRequest -Uri $downloadUrl -Out $pathZip
    Remove-Item -Path $pathExtract -Recurse -Force -ErrorAction SilentlyContinue
    Expand-Archive -Path $pathZip -DestinationPath $pathExtract -Force
    Remove-Item $pathZip -Force
    #    Move-Item "$pathExtract\booty.exe" "$pathExtract\booty"
}

function add_envpath
{
    param(
        [Parameter(Mandatory = $true)]
        [string] $Path,

        [ValidateSet('Machine', 'User', 'Session')]
        [string] $Container = 'Session'
    )

    if ($Container -ne 'Session')
    {
        $containerMapping = @{
            Machine = [EnvironmentVariableTarget]::Machine
            User = [EnvironmentVariableTarget]::User
        }
        $containerType = $containerMapping[$Container]

        $persistedPaths = [Environment]::GetEnvironmentVariable('Path', $containerType) -split ';'
        if ($persistedPaths -notcontains $Path)
        {
            $persistedPaths = $persistedPaths + $Path | Where-Object { $_ }
            [Environment]::SetEnvironmentVariable('Path', $persistedPaths -join ';', $containerType)
        }
    }

    $envPaths = $env:Path -split ';'
    if ($envPaths -notcontains $Path)
    {
        $envPaths = $envPaths + $Path | Where-Object { $_ }
        $env:Path = $envPaths -join ';'
    }
}

download_latest
add_envpath $pathExtract "User"

Write-Host 'For more information, see: https://github.com/amplify-edge/booty'
