Name:           quadlet-compose-generator
Version:        0.1.0
Release:        1%{?dist}
Summary:        Systemd user generator for Podman Compose files
License:        MIT
URL:            https://github.com/kuyacarlo/quadlet-compose-generator
Source0:        %{name}-%{version}.tar.gz
BuildRequires:  golang >= 1.22
BuildRequires:  git
Provides:       bundled(golang(*))

%description
quadlet-compose-generator is a systemd user generator that auto-discovers Podman Compose
files in ~/.config/containers/compose/ and generates systemd service units on
daemon-reload. Drop a compose.yml, reload, start.

%prep
%autosetup -n %{name}-%{version}

%build
export GOFLAGS="-mod=vendor"
export GOTOOLCHAIN=local
go build -o %{name} -ldflags "-s -w" .

%install
install -Dpm 0755 %{name} %{buildroot}%{_prefix}/lib/systemd/user-generators/%{name}

%check
export GOFLAGS="-mod=vendor"
export GOTOOLCHAIN=local
go test ./...

%files
%license LICENSE
%{_prefix}/lib/systemd/user-generators/%{name}

%changelog
* Sat Aug 02 2026 kuyacarlo <kuyacarlo@users.noreply.github.com> - 0.1.0-1
- Initial package
