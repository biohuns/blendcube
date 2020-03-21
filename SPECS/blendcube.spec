%define debug_package %{nil}

Name: blendcube
Version: %{_version}
Release: %{_release}%{?_dist}

Summary: Rubik's Cube 3D Model Server
License: MIT
Group: Applications/Internet
URL: https://github.com/biohuns/blendcube

Source0: %{name}-%{version}.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
Simple API Server for Generating Rubik's Cube 3D Model from URL

%prep
%setup -q

%build

%install
install -D %{name} %{buildroot}%{_bindir}/%{name}
install -D config.json %{buildroot}%{_sysconfdir}/%{name}/config.json
install -D cube.gltf %{buildroot}%{_sysconfdir}/%{name}/cube.gltf
install -D cube.glb %{buildroot}%{_sysconfdir}/%{name}/cube.glb
install -D service %{buildroot}/etc/systemd/system/%{name}.service
install -D logrotate %{buildroot}/etc/logrotate.d/%{name}
install -d %{buildroot}/var/log/%{name}

%files
%attr(0755,root,root) %{_bindir}/%{name}
%defattr(0644,root,root, 0755)
%config(noreplace) %{_sysconfdir}/%{name}/config.json
%config(noreplace) %{_sysconfdir}/%{name}/cube.gltf
%config(noreplace) %{_sysconfdir}/%{name}/cube.glb
%config(noreplace) %{_sysconfdir}/systemd/system/%{name}.service
%config(noreplace) %{_sysconfdir}/logrotate.d/%{name}
/var/log/%{name}

%post
if [ $1 -eq 1 ]; then
    /bin/systemctl enable %{name}.service
fi
/bin/systemctl stop %{name}.service || :
/bin/systemctl daemon-reload >/dev/null 2>&1 || :
/bin/systemctl start %{name}.service

%preun
if [ $1 -eq 0 ]; then
    /bin/systemctl stop %{name}.service || :
fi

%clean
rm -rf %{buildroot}
