const routes = [
    { path: '/', component: sitelistTemplate },
    { path: '/sitelist', component: sitelistTemplate },
    { path: '/addsite', component: addsiteTemplate },
    { path: '/editCaddySiteListConfig', name: 'editCaddySiteListConfig', component: editCaddySiteListConfigTemplate },
    { path: '/loglist/:ProjectID', name: 'loglist', component: loglistTemplate },
    { path: '/logdetail/:TraceID', name: 'logdetail', component: logDetailTemplate },
    { path: '/jsonconfig', name: 'jsonconfig', component: jsonConfigTemplate }
]

const router = new VueRouter({
    routes // short for `routes: routes`
})