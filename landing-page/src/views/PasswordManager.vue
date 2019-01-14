<template>
  <div class="row" id="main-content">
    <div class="col-md-10">
        <h1 class="super-header">MVP</h1>
        <p class="subtitle">Minimum Remarkable (Viable) Product</p>
        <h3 class="title">Introduction</h3>
        <p>The ambitious goals of Dynauth is to provide seamless, secure, and self-owned identity to everyone. However, this is a monumental challenge that requires loads of tech to make everything function holistically. So where do we start? In a word: password manager.</p>
        <h3 class="title">A password manager?</h3>
        <p class="mb-5"><em>Hasn't that been done before?</em></p>
        <p>Well yes, it has. In many different forms. By many different organizations. And that's great! Password managers are a large improvement over passwords, the only problem is that they still rely on passwords as the final barrier to unwanted entry.</p>
        <p>Since the goal of Dynauth is to provide <em>secure</em> and <em>seamless</em> identification, that forces us to:
        <ol type="a">
            <li>Get rid of passwords as an authentication source since they are no longer secure... and</li>
            <li>Seamlessly authenticate users into legacy applications that require... passwords</li>
        </ol>
        Well that's a bit of a paradox isn't it? We want to get rid of passwords... but we're forced to accomodate them.</p>
        <p>Unfortunately, just because there is a better alternative to passwords doesn't mean we can move on immediately. Changes like that take time -- lots of time. We must accomodate the present to develop for the future, and the present is passwords.</p>
        <h3 class="title">So what's the MVP exactly?</h3>
        <p>The MVP will be a simple password manager that implements the technology behind dynamic authentication to create a more secure alternative to anything currently out there.</p>
        <h3 class="title">Features</h3>
        <p>There are some improvements that not having a master password will provide.</p>
        <h3 class="title">Current Progress</h3>
        <p>The first iteration of the Dynauth Password Manager (name pending) is available here: <a href="https://github.com/mainsailstudio/dynauth-password-manager">github.com/mainsailstudio/dynauth-password-manager</a></p>
        <p>It doesn't really do much yet.</p>
    </div>
  </div>
</template>
<script>
export default {
  data: ()=> ({

  }),
  mounted(){
 
  },
  methods: {

  },
}
</script>