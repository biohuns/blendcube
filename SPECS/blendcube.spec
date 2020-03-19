Name: blendcube
Version: 0.4.0
Release: 1
Summary: Simple API Server for Generating Rubik's Cube 3D Model from URL
License: MIT
Group: Development/Tools
URL: https://github.com/biohuns/blendcube
Source: https://codeload.github.com/biohuns/blendcube/tar.gz/v0.3.0
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root
BuildRequires: make

%description
Simple API Server for Generating Rubik's Cube 3D Model from URL

%prep
rm -rf %{buildroot}

%build

%install
ls -la /github/home/rpmbuild/
ls -la /github/home/rpmbuild/BUILD/
install -D %{name} %{buildroot}%{_bindir}/%{name}
install -D config.json.example %{buildroot}%{_sysconfdir}/%{name}/config.json
install -D model/* %{buildroot}%{_sysconfdir}/%{name}/model/
install -d %{buildroot}/var/log/%{name}

%files
%defattr(0755,root,root)
%{_bindir}/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/config.json
%config(noreplace) %{_sysconfdir}/%{name}/cube.gltf
%config(noreplace) %{_sysconfdir}/%{name}/cube.glb
/var/log/%{name}

%clean
rm -rf %{buildroot}
