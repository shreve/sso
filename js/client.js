let SSOperhero = (() => {
  let log = console.log.bind(window.console, '[SSOperhero Client]')
  if (location.hostname !== 'localhost') { log = () => {} }

  let opts = {
    login_success: () => { log('Logged in successfully') },
    logout_success: () => { log('Logged out successfully') },
    login_expired: () => { log('Login has expired') }
  }

  let extend = (orig, add) => {
    for (var prop in add) {
      if (add.hasOwnProperty(prop)) {
        orig[prop] = add[prop]
      }
    }
    return orig
  }

  let token = () => {
    let t = localStorage.getItem('token')
    if (t) {
      log('[token]', t)

      let bits = t.split('.')

      return {
        header: JSON.parse(atob(bits[0])),
        ...JSON.parse(atob(bits[1]))
      }
    }
  }

  let post_message = (message) => {
    log('[message out]', message)
    provider.contentWindow.postMessage(message, opts.provider)
  }

  let receive_message = (event) => {
    if (event.origin != opts.provider) return;
    log('[message in]', event)
    switch (event.data.intent) {
    case 'token:set':
      logged_in = localStorage.getItem('token')
      localStorage.setItem('token', event.data.value)
      if (!logged_in) opts.login_success()
      break;
    case 'token:clear':
      logged_in = localStorage.getItem('token')
      localStorage.removeItem('token')
      if (logged_in) opts.logout_success()
      break;
    case 'error':
      break;
    }
  }

  let add_provider_window = () => {
    if (document.getElementById('ssoperhero-provider')) { return }

    provider = document.createElement('iframe')
    provider.id = 'ssoperhero-provider'
    provider.src = opts.provider
    provider.style.width = provider.style.height = provider.style.border = '0'
    document.body.appendChild(provider)
  }

  return {
    init: (options) => {
      opts = extend(opts, options)

      window.addEventListener('message', receive_message, false)

      add_provider_window();
    },

    logged_in: () => {
      let t = token()
      return !!t && (t.exp * 1000) > (new Date())
    },

    login: (args) => {
      post_message({ intent: 'login',
                     username: args.username,
                     password: args.password })
    },

    logout: () => {
      post_message({ intent: 'logout' })
    },

    register: (args) => {
      post_message({ intent: 'register',
                     username: args.username,
                     password: args.password })
    },

    token: token
  }
})()

url = new URL(document.getElementById('ssoperhero').src)
SSOperhero.init({ provider: url.origin })
