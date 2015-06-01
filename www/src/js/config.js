module.exports = {
  PACKAGES: '/api/v2/packages/',
  SETTINGS: '/api/v2/packages/ubuntu-core',
  FILTERED_SNAPS: [
    'webdm',
    'ubuntu-core'
  ],
  NON_INSTALLABLE_IDS: [
    'webdm',
    'ubuntu-core.ubuntu'
  ],
  NON_INSTALLABLE_TYPES: [
    'oem'
  ],
  INSTALL_STATE: {
    INSTALLED: 'installed',
    INSTALLING: 'installing',
    REMOVED: 'removeed',
    REMOVING: 'removeing'
  },
  INSTALL_POLL_WAIT: 16 // milliseconds
};
