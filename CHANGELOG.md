# Changelog

## 2022-06-22 - Extend API

### Moved
- API url to .env (SRC_API parameter)

### Added
- "Status" field to outages retrieved from database. This shows whether an outage is still active or not
- after_end_date and before_start_date filter parameters added

## 2022/06/02 - Fix for #11

### Added
- Yearly cronjob to edit poorly saved addresses and related API

### Fixed
- Short-hands in street and suburb names
- Street numbers, postcodes and extra comma divisors from being recognised
- Inconsistent capitalisation of addresses