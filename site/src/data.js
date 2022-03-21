const columns = [
  { title: "Provider", field: "name", filtering:false, render: rowData => <a href={rowData.url}>{rowData.name}</a> },
  { title: "", field: "docs", filtering:false, render: rowData => <a href={rowData.docsUrl}>Docs</a> },
  { title: "Updated", field: "updated", filtering:false},
  {
    title: "CRDs maturity", field: "crdsMaturity",
    lookup: { Unreleased: "Unreleased", Alpha: "Alpha", Beta: "Beta", V1: "V1" },
    defaultFilter: ["Alpha", "Beta", "V1"]
  },
  { title: "CRDs", field: "crds", filtering:false, type: "numeric" },
];


const data = [
  { name: "crossplane-contrib/provider-jet-gcp", url: "https://github.com/crossplane-contrib/provider-jet-gcp", docsUrl: "https://github.com/crossplane-contrib/provider-jet-gcp", updated: "2022-03-19", crdsMaturity: "Unreleased", crds: 1995 },
  { name: "crossplane-contrib/provider-jet-gcp", url: "https://github.com/crossplane-contrib/provider-jet-gcp", docsUrl: "https://github.com/crossplane-contrib/provider-jet-gcp", updated: "2022-03-19", crdsMaturity: "Alpha", crds: 1995 },
  { name: "crossplane-contrib/provider-jet-gcp", url: "https://github.com/crossplane-contrib/provider-jet-gcp", docsUrl: "https://github.com/crossplane-contrib/provider-jet-gcp", updated: "2022-03-19", crdsMaturity: "Alpha", crds: 1995 },
  { name: "crossplane-contrib/provider-jet-gcp", url: "https://github.com/crossplane-contrib/provider-jet-gcp", docsUrl: "https://github.com/crossplane-contrib/provider-jet-gcp", updated: "2022-03-19", crdsMaturity: "Alpha", crds: 1995 },
];

const exported = { data: data, columns: columns }

export default exported;
