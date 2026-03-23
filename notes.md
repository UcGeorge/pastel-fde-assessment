## Notes taken by George Uche-Umeh during completion of the assessment

- The API returns 200 with body {"message":"Access Denied"}
- The provided documentation has a typo "sacntion" in Check Sanction (Instant) [POST api/v1/aml/sacntion/instant]. Using this typo returns a 404 and thereby blocks any IP address making that request. One would have to manually fix the typo for requests to go through.
- The check PEP and CheckSanctions endpoints both return count of zero, whereas there are elements or objects within the list in the data field.
- The URLs in the datasets and some sanction responses and PEP responses are empty strings.