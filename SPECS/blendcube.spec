%define name        blendcube
%define version     0.4.0
%define release     1
%define buildroot   %{_topdir}

Name: %{name}
Version: %{version}
Release: %{release}
Summary: Simple API Server for Generating Rubik's Cube 3D Model from URL
License: MIT
Group: Development/Tools
URL: https://github.com/biohuns/blendcube
Source: https://codeload.github.com/biohuns/blendcube/tar.gz/v0.3.0
BuildRoot: %{buildroot}
BuildRequires: make go

%description
Simple API Server for Generating Rubik's Cube 3D Model from URL

%prep
rm -rf %{buildroot}

%build
make

%install
cp %{name} %{buildroot}/usr/bin/%{name}
cp config.json %{buildroot}/etc/%{name}/config.json

%files
%defattr(0755,root,root)
/usr/bin/%{name}

%config
%config(noreplace) /etc/%{name}/config.json
%config(noreplace) /etc/%{name}/cube.gltf
%config(noreplace) /etc/%{name}/cube.glb
