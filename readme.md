# Pipeline Database API 

longterm goal is to leverage that API and build a UI where users could things like:
- see what versions of RKE2 are in the secured registry
- validate the images in the registry are signed and when they were lasted trivy scanned
- download a pre-compiled tarball for a specific release of a product

## REST Schema
>prefix: https://\<backendurl\>/api/v0/)

core api:
- /product
    - GET: get all products
    - POST: create new product
- /product/{product_name}
    - GET: get product
    - PUT: update product
    - DELETE: delete product
- /product/{product_name}/release
    - GET: get all releases for product
    - POST: create new release for product
- /product/{product_name}/release/{release_name}
    - GET: get release
    - PUT: update release
    - DELETE: delete release
- /product/{product_name}/release/{release_name}/image
    - GET: get all images for product release
    - POST: create new image for product release
- /product/{product_name}/release/{release_name}/image/{image_name}
    - GET: get image from product release
    - PUT: update image from product release
    - DELETE: delete image from product release

other:
- /releaseimgmapping
    - GET: get all images for product release
    - POST: create new image for product release
- /releaseimgmapping/{releaseimgmapping_id}
    - GET: get releaseimgmapping
    - PUT: update releaseimgmapping
    - DELETE: delete releaseimgmapping
- /user
    - POST: should accept username and password
        - only returns cookie first time user is created
- /login
    - POST: should accept username and password
        - returns authentication cookie/token (stored in browser by frontend)
