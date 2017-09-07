# Template RPM Spec file. Please update all references to goms with your
# service name. The package name must always start with `yapo-` or $god will
# kill a pretty kitten and deliver its pieces to you
#

Name:       yapo-goms
Version:    0.00.1
Release:    %(expr `date +%s`)
Summary:    Service template in go

License:    Copyright Schibsted Classified Media 2017
URL:        https://github.schibsted.io/Yapo/goms
Source:     https://github.schibsted.io/Yapo/goms

BuildRequires:	golang >= 1.3

%define svc goms
%define _topdir %(pwd)/build
%define codesrc %(pwd)
%define bindir /opt/%{svc}
%define confdir /opt/%{svc}/conf
%define initdir /etc/init.d
%define _sysconfdir /etc/sysconfig

%description
Micro Service template to serve JSON APIs in Golang

%pre
if ! id %{svc} &>/dev/null; then useradd %{svc} ; fi
%(sed -i 's/__API_NAME__/%{svc}/g' %{codesrc}/scripts/api)

%post
%(sed 's/__API_NAME__/%{svc}/g' scripts/post-install.sh)

%install
install -d %{buildroot}%{bindir}
install -d %{buildroot}%{confdir}
install -d %{buildroot}%{initdir}
install -p -m 0755 %{codesrc}/%{svc} %{buildroot}%{bindir}/%{svc}-api
install -p -m 0755 %{codesrc}/scripts/api %{buildroot}%{initdir}/%{svc}-api
install -p -m 0644 %{codesrc}/conf/conf.json.in %{buildroot}%{confdir}/conf.json

%clean
rm -rf %{buildroot}

%files
%defattr(-,%{svc},%{svc},-)
%{bindir}/%{svc}-api
%{confdir}/conf.json
%attr(-, root, root) %{initdir}/%{svc}-api
