const routes = [
    { path: '/', component: sitelistTemplate },
    { path: '/sitelist', component: sitelistTemplate },
    { path: '/addsite', component: addsiteTemplate },
    { path: '/editCaddySiteListConfig', name: 'editCaddySiteListConfig', component: editCaddySiteListConfigTemplate },
    { path: '/caddyConfig', name: 'caddyConfig', component: caddyConfigTemplate },
    { path: '/pwdmanage', name: 'pwdmanage', component: pwdmanageTemplate },
    { path: '/certlist', name: 'certlist', component: certlistTemplate },
    { path: '/webapplist', name: 'webapplist', component: webapplistTemplate },
]

const router = new VueRouter({
    routes // short for `routes: routes`
})