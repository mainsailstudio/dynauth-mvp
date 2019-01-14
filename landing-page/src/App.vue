<template>
  <div id="app">
    <div class="container-fluid">
        <div class="row">
          <nav class="col-md-2 d-md-block sidebar">
            <div class="sidebar-sticky" id="sidebar-content">
                <button class="btn btn-link bd-search-docs-toggle d-md-none p-0 ml-3 collapsed" type="button" data-toggle="collapse" data-target="#sidebar-nav" aria-controls="sidebar-nav" aria-expanded="false" aria-label="Toggle sidebar navigation"><svg style="margin-top:6px;" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 30 30" width="30" height="30" focusable="false"><title>Menu</title><path stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-miterlimit="10" d="M4 7h22M4 15h22M4 23h22"></path></svg>
                </button>
              <ul class="nav collapse" id="sidebar-nav">
                <li class="nav-item">
                  <router-link to="/" class="nav-link">
                    Home
                  </router-link>
                </li>
                <li class="nav-item">
                  <router-link to="/dynamic-authentication" class="nav-link">
                    Dyamic Authentication
                  </router-link>
                </li>
                <li class="nav-item">
                  <router-link to="/password-manager" class="nav-link">
                    Password Manager (MVP)
                  </router-link>
                </li>
                <!-- <li class="nav-item">
                  <router-link to="/about" class="nav-link">
                    About
                  </router-link>
                </li> -->
                <li class="nav-item">
                  <router-link to="/contact" class="nav-link">
                    Contact
                  </router-link>
                </li>
              </ul>
            </div>
          </nav>
          <main role="main" class="col-md-9 ml-sm-auto col-lg-10 px-4">
              <div id="main-content">
                <router-view/>
              </div>
          </main>
      </div>
    </div>
  </div>
</template>
<style>
  @import './assets/css/resume.css';
  @import './assets/vendor/fontawesome-free-5.6.3/css/all.min.css';
</style>