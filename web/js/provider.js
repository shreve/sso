let SSO = (() => {
  let opts = {}
  let log = console.log.bind(window.console, '[SSO Provider]')

  let request = (opt) => {
    log('[request]', opt)
    if (!opt.method) opt.method = 'GET'

    let req = new XMLHttpRequest()
    req.open(opt.method, opt.url, true)
    req.setRequestHeader('X-Origin', document.referrer)
    req.onload = () => {
      if (req.responseText === "")
        data = {}
      else
        data = JSON.parse(req.responseText)

      if (200 <= req.status && req.status < 400)
        opt.success(data)
      else
        opt.error(data)
    }
    req.onerror = () => {
      post_message({ intent: 'error', value: req })
    }

    if (opt.method === 'GET')
      req.send()
    else {
      req.setRequestHeader('Content-Type', 'application/json;charset=utf-8')
      req.send(JSON.stringify(opt.data))
    }
  }

  //
  // The Remote Requests
  //
  // These 4 functions talk to the server
  // Fetch a token, register account, log in, log out
  let request_token = () => {
    request({
      url: '/status',
      success: send_token,
      error: () => {}
    })
  }

  let register = (args) => {
    request({
      method: 'POST',
      url: '/register',
      data: args,
      success: send_token,
      error: send_error
    })
  }

  let login = (args) => {
    request({
      method: 'POST',
      url: '/login',
      data: args,
      success: send_token,
      error: send_error
    })
  }

  let logout = () => {
    request({
      url: '/logout',
      success: (data) => {
        post_message({ intent: 'token:clear' })
      },
      error: send_error
    })
  }


  //
  // The Message Protocol
  //
  // Post message and receive message handle talking to the client through the
  // window event.
  // Send methods are helpers for request handlers.
  let post_message = (message) => {
    if (!opts.client) { return }
    log('[message out]', message)
    log('[opts]', opts)
    window.parent.postMessage(message, window.location.protocol + '//' + opts.client)
  }

  let receive_message = (event) => {
    log('[message in]', event);
    switch (event.data.intent) {
    case 'token:get':
      request_token();
      break;
    case 'logout':
      logout();
      break;
    case 'login':
      login(event.data);
      break;
    case 'register':
      register(event.data);
      break;
    }
  }

  let send_token = (data) => {
    post_message({ intent: 'token:set', value: data.token })
  }

  let send_error = (data) => {
    post_message({ intent: 'error', value: data.error })
  }

  // Public API
  return {
    init: options => {
      opts = options
      window.addEventListener('message', receive_message, false)
      request_token();
    }
  }
})();
