# GPT-RAY

Although I have a solid 10 years programming experience, I have no prior knowledge of light transport and rendering software. 

So I thought why not let ChatGPT help me write a raytracer.

As I am in no position to judge the correctness of the code (or at least I will pretend not to be, for the sake of the experiment), I will have to try to keep improving iteratively through further chat interaction.

#### 01. Initial render

The basis of the code, up until the point it could output an image for the first time, was generated within one and a half hours!
This is the initial render, personally I think this is an amazing result in such a small timeframe, although there are probably errors in the code.

![Render 01](renders/01.png)

#### 02. Clamping

Adding clamping clearly improves the output image.

![Render 02](renders/02.png)

#### 03. Adding a Plane

Adding a plane to the scene seems to create a proper mess. Time for some chat guided debugging...

![Render 03](renders/03.png)