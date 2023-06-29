Group: BytemanD
Name: ec-tools
Version: VERSION
Release: 1
Summary: Ec Tools
License: ASL 2.0

Source0: ec-tools
Source1: ec-tools-template.yaml

Requires: libvirt-devel

%global CONFIG_DIRNAME ec-tools
%global CONFIG_PATH /etc/${CONFIG_DIRNAME}

%description
Golang EC Tools


%prep
#cp -p %SOURCE0 %{_builddir}
mkdir -p %{_builddir}${CONFIG_PATH}


%files
%{_bindir}/ec-tools
%{_sysconfdir}/ec-tools/ec-tools-template.yaml

%install
install -m 755 -d %{buildroot}%{_bindir}
install -m 755 -d %{buildroot}%{_sysconfdir}/%{CONFIG_DIRNAME}

install -p -m 755 -t %{buildroot}%{_bindir} %{SOURCE0}
install -p -m 755 -t %{buildroot}%{_sysconfdir}/%{CONFIG_DIRNAME} %{SOURCE1}

%post

cd %{_sysconfdir}/ec-tools/
if [[ ! -f ec-tools.yaml ]]; then
    cp ec-tools-template.yaml ec-tools.yaml
fi
