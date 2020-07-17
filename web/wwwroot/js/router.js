const routes = [
    { path: '/', component: sitelistTemplate },
    { path: '/sitelist', component: sitelistTemplate },
    { path: '/addsite', component: addsiteTemplate },
    { path: '/editCaddySiteListConfig', name: 'editCaddySiteListConfig', component: editCaddySiteListConfigTemplate },
    { path: '/jsonconfig', name: 'jsonconfig', component: jsonConfigTemplate },
    { path: '/pwdmanage', name: 'pwdmanage', component: pwdmanageTemplate },
    { path: '/certlist', name: 'certlist', component: certlistTemplate },
]

const router = new VueRouter({
    routes // short for `routes: routes`
})