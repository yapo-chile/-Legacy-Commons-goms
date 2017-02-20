Name:		goms
Version:	0.00.0
Release:	%(expr `date +%s`)
Summary:	Service in go

License:	Copyright Schibsted Classified Media 2017
URL:		https://github.schibsted.io/Yapo/goms
Source:		https://github.schibsted.io/Yapo/goms

BuildRequires:	golang >= 1.3

%define _topdir %(pwd)/build
%define codesrc %(pwd)
%define mspkg src/github.schibsted.io/Yapo/goms
%define bindir /opt/goms
%define confdir /opt/goms/conf
%define initdir /etc/init.d
%define _sysconfdir /etc/sysconfig

%description
MS escrito en Goland como una API JSON.

%pre
if ! id goms &>/dev/null; then useradd goms ; fi

%post
%(cat scripts/goms-post-install.sh)

%install
install -d %{buildroot}%{bindir}
install -d %{buildroot}%{confdir}
install -d %{buildroot}%{initdir}
install -p -m 0755 %{codesrc}/goms%{buildroot}%{bindir}/goms-api
install -p -m 0755 %{codesrc}/scripts/goms-api %{buildroot}%{initdir}
install -p -m 0644 %{codesrc}/conf/conf.json.in %{buildroot}%{confdir}/conf.json

%clean
rm -rf %{buildroot}

%files
%defattr(-,goms,goms,-)
%{bindir}/goms-api
%{confdir}/conf.json
%attr(-, root, root) %{initdir}/goms-api
